.PHONY: help setup run down logs test lint clean migrate

help:
	@echo "Revenant Backend - Available Commands"
	@echo "====================================="
	@echo "make setup      - Install dependencies"
	@echo "make run        - Start the application with docker-compose"
	@echo "make down       - Stop the application"
	@echo "make logs       - Show application logs"
	@echo "make test       - Run tests"
	@echo "make lint       - Run linter"
	@echo "make clean      - Clean build artifacts"

setup:
	go mod download
	go mod tidy

run:
	docker-compose up --build

down:
	docker-compose down

logs:
	docker-compose logs -f

test:
	go test -v -cover ./...

lint:
	go vet ./...

clean:
	rm -f revenant-app
	go clean
