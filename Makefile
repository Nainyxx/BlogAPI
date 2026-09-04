.PHONY: help install run build test clean fmt vet lint docker-build docker-up docker-down db-migrate

# Переменные
GO := go
BINARY_NAME := blog-api
MAIN_FILE := main.go
PORT := 8082

help:
	@echo "Blog API - Available Commands:"
	@echo ""
	@echo "  make install      - Download and install dependencies"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the binary"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format the code"
	@echo "  make vet          - Run go vet"
	@echo "  make lint         - Run linter (golangci-lint)"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make db-migrate   - Apply db_schema.sql to the database from .env"
	@echo "  make help         - Show this help message"

install:
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependencies installed"

run:
	$(GO) run $(MAIN_FILE)

build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Built: $(BINARY_NAME)"

test:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

fmt:
	$(GO) fmt ./...
	@echo "✅ Code formatted"

vet:
	$(GO) vet ./...
	@echo "✅ Go vet passed"

lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "✅ Lint passed"

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	@echo "✅ Cleaned"

docker-build:
	docker build -t $(BINARY_NAME):latest .
	@echo "✅ Docker image built"

docker-up:
	docker-compose up -d
	@echo "✅ Docker containers started"
	@echo "   API:     http://localhost:8080"
	@echo "   pgAdmin: http://localhost:5050"

docker-down:
	docker-compose down
	@echo "✅ Docker containers stopped"

docker-logs:
	docker-compose logs -f blog_api

db-migrate:
	@set -a && . ./.env && set +a && \
	PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d $$DB_NAME -f db_schema.sql
	@echo "✅ Schema applied"

# Development targets
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

deps-update:
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "✅ Dependencies updated"

# Quick start
quick-start: install run

# Full setup
setup: install fmt vet test build
	@echo "✅ Setup complete"
