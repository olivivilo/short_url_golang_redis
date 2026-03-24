.PHONY: help build run test test-unit test-integration test-coverage clean docker-build docker-up docker-down fmt vet lint

# Default target
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests only"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  clean              - Clean build artifacts"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-up          - Start services with Docker Compose"
	@echo "  docker-down        - Stop services with Docker Compose"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  lint               - Run golangci-lint (if installed)"

# Build the application
build:
	@echo "Building..."
	@go build -v -o bin/shorturl ./cmd/shorturl

# Run the application
run:
	@echo "Running..."
	@go run ./cmd/shorturl/main.go

# Run all tests
test:
	@echo "Running all tests..."
	@go test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v ./internal/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -v ./test/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t short-url-service:latest -f deploy/Dockerfile .

# Start services with Docker Compose
docker-up:
	@echo "Starting services..."
	@docker-compose -f deploy/docker-compose.yml up -d

# Stop services with Docker Compose
docker-down:
	@echo "Stopping services..."
	@docker-compose -f deploy/docker-compose.yml down

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run golangci-lint (if installed)
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
