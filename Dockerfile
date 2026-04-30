FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o revenant-app main.go

FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates postgresql-client

COPY --from=builder /app/revenant-app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./revenant-app"]
