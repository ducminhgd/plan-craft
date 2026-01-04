.PHONY: help build run clean test deps migrate db-clean

# Default target
help:
	@echo "Available targets:"
	@echo "  deps      - Download and install dependencies"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  db-clean  - Remove database file"
	@echo "  migrate   - Run database migrations"

# Download and install dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build the application
build: deps
	@echo "Building application..."
	go build -o bin/plancraft cmd/server/main.go

# Run the application
run: build
	@echo "Running application..."
	./bin/plancraft

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f plancraft

# Remove database file
db-clean:
	@echo "Removing database file..."
	rm -f data/plancraft.db
	rm -f data/plancraft.db-shm
	rm -f data/plancraft.db-wal

# Run database migrations (same as running the app, which auto-migrates)
migrate: build
	@echo "Running database migrations..."
	./bin/plancraft

# Development mode - run with auto-reload (requires air)
dev:
	@echo "Starting development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Install it with: go install github.com/air-verse/air@latest"; \
		echo "Running without auto-reload..."; \
		go run cmd/server/main.go; \
	fi

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

