.PHONY: build run test clean dev help

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  dev      - Run the application with auto-reload"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  help     - Show this help message"

# Build the application
build:
	@echo "Building the application..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build completed. Binary: bin/server"

# Run the application
run:
	@echo "Running the application..."
	@go run cmd/server/main.go

# Run with development settings
dev:
	@echo "Running in development mode..."
	@air -c .air.toml || go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f ecommerce.db
	@rm -f coverage.out coverage.html
	@echo "Clean completed"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"

# Database operations
db-reset:
	@echo "Resetting database..."
	@rm -f ecommerce.db
	@echo "Database reset completed"

# Docker operations
docker-build:
	@echo "Building Docker image..."
	@docker build -t Product Order System:latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --rm Product Order System:latest

# Development setup
setup: deps
	@echo "Setting up development environment..."
	@cp .env.example .env 2>/dev/null || touch .env
	@echo "Setup completed. Edit .env file if needed."