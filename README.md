# revenant-backend

REST API for Revenant — personal and shared expense tracking. Built with Go 1.23, Gin, PostgreSQL, and JWT auth.

## Stack

- **Go 1.23** · **Gin** · **PostgreSQL 15** · **Docker**
- `golang-migrate` — schema migrations
- `zap` — structured logging
- `validator` — request validation
- `JWT (HS256)` — stateless auth

## Quick Start

```bash
# With Docker (recommended)
make run        # docker-compose up --build
make logs       # tail logs
make down       # stop

# Local dev (requires Postgres running separately)
make setup      # go mod download + tidy
go run main.go
```

## Environment Variables

See `.env.example`. Key vars:

```env
SERVER_PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=revenant
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_SSLMODE=disable
JWT_SECRET=change-this-in-production
JWT_EXPIRATION=24h
LOG_LEVEL=info
```

## API Reference

All `/api/v1/*` endpoints except auth require `Authorization: Bearer <token>`.

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/auth/register` | Register (email, password, first_name, last_name, username) |
| POST | `/api/v1/auth/login` | Login → returns JWT token |

### Users
| Method | Path | Description |
|---|---|---|
| GET | `/api/v1/users/search?username=` | Search user by username (for expense sharing) |
| GET | `/api/v1/users/:id` | Get user by UUID |

### Expenses
| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/expenses` | Create expense |
| GET | `/api/v1/expenses?type=personal\|shared` | List expenses |
| GET | `/api/v1/expenses/:id` | Get expense with participants |
| PUT | `/api/v1/expenses/:id` | Update expense (creator only) |
| DELETE | `/api/v1/expenses/:id` | Delete expense (creator only) |
| POST | `/api/v1/expenses/:id/participants` | Add participant (creator only) |

**Currencies:** CLP (default), CLF, USD, EUR

**Split types:** `personal`, `equal`, `percentage`, `exact`

### Create Expense — request body

```json
{
  "title": "Coffee",
  "amount": 5000,
  "currency": "CLP",
  "category": "Dining",
  "date": "2026-05-12",
  "description": "",
  "split_type": "personal",
  "participants": []
}
```

## Project Structure

```
internal/
├── config/          # App config (env vars)
├── database/        # DB connection + migrations runner
├── logger/          # Zap logger setup
├── models/          # Data structs (User, Expense, etc.)
├── services/        # Business logic
└── server/
    ├── handlers/    # HTTP handlers
    └── middleware/  # JWT auth middleware
migrations/          # SQL migration files (golang-migrate)
```

## Useful Commands

```bash
make setup    # download dependencies
make run      # start with Docker
make test     # go test -v -cover ./...
make lint     # go vet ./...
make down     # stop Docker services
make clean    # clean build artifacts
```

## Adding a Migration

```bash
migrate create -ext sql -dir migrations -seq create_<table>_table
# Edit the generated .up.sql and .down.sql files
```
