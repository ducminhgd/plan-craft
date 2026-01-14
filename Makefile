.PHONY: help build run clean test deps migrate migrate-up migrate-down migrate-status db-clean frontend-install frontend-build frontend-dev wails-dev wails-build

# Default target
help:
	@echo "Available targets:"
	@echo "  deps             - Download and install Go dependencies"
	@echo "  frontend-install - Install frontend dependencies"
	@echo "  frontend-build   - Build frontend for production"
	@echo "  frontend-dev     - Run frontend dev server"
	@echo "  wails-dev        - Run Wails in development mode (recommended)"
	@echo "  wails-build      - Build Wails application for production"
	@echo "  build            - Build the Wails application"
	@echo "  run              - Run the Wails application in dev mode"
	@echo "  test             - Run tests"
	@echo "  clean            - Clean build artifacts"
	@echo "  db-clean         - Remove database file"
	@echo "  migrate          - Run all pending database migrations (up)"
	@echo "  migrate-up       - Run all pending database migrations"
	@echo "  migrate-down     - Rollback the last migration"
	@echo "  migrate-status   - Show current migration status"

# Download and install Go dependencies
deps:
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy

# Install frontend dependencies
frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Build frontend for production
frontend-build: frontend-install
	@echo "Building frontend..."
	cd frontend && npm run build

# Run frontend dev server
frontend-dev:
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

# Run Wails in development mode (recommended)
wails-dev: frontend-install
	@echo "Starting Wails development mode..."
	@if command -v wails > /dev/null; then \
		wails dev; \
	else \
		echo "❌ Wails is not installed. Install it with:"; \
		echo "   go install github.com/wailsapp/wails/v2/cmd/wails@latest"; \
		exit 1; \
	fi

# Build Wails application for production
wails-build: deps frontend-build
	@echo "Building Wails application..."
	@if command -v wails > /dev/null; then \
		wails build; \
	else \
		echo "❌ Wails is not installed. Install it with:"; \
		echo "   go install github.com/wailsapp/wails/v2/cmd/wails@latest"; \
		exit 1; \
	fi

# Build the Wails application
build: wails-build

# Run the Wails application in dev mode
run: wails-dev

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf build/bin/
	rm -f plancraft
	rm -rf frontend/dist/*
	@echo "✅ Clean complete"

# Remove database file
db-clean:
	@echo "Removing database file..."
	rm -f data/plancraft.db
	rm -f data/plancraft.db-shm
	rm -f data/plancraft.db-wal

# Database configuration
DB_DSN ?= data/plancraft.db
MIGRATE_CMD = migrate -path migrations -database "sqlite3://$(DB_DSN)"

# Run all pending database migrations (up)
migrate: migrate-up

# Run all pending database migrations
migrate-up:
	@echo "Running database migrations..."
	@mkdir -p data
	@if command -v migrate > /dev/null; then \
		$(MIGRATE_CMD) up; \
	else \
		echo "❌ golang-migrate is not installed. Install it with:"; \
		echo "   go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	@if command -v migrate > /dev/null; then \
		$(MIGRATE_CMD) down 1; \
	else \
		echo "❌ golang-migrate is not installed. Install it with:"; \
		echo "   go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Show current migration status
migrate-status:
	@echo "Checking migration status..."
	@if command -v migrate > /dev/null; then \
		$(MIGRATE_CMD) version; \
	else \
		echo "❌ golang-migrate is not installed. Install it with:"; \
		echo "   go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Development mode - run Wails dev (alias for run)
dev: wails-dev

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
	fi

# Generate mocks (if using mockery)
mocks:
	@echo "Generating mocks..."
	@if command -v mockery > /dev/null; then \
		mockery --all --dir internal --output internal/mocks; \
	else \
		echo "mockery not found. Install it with: go install github.com/vektra/mockery/v2@latest"; \
	fi

