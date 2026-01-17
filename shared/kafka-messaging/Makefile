# Kafka Messaging Package Makefile

.PHONY: test test-unit test-integration build clean lint

# Default target
all: test build

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v ./...

# Run integration tests (requires running Kafka)
test-integration:
	@echo "Running integration tests..."
	@echo "Note: This requires a running Kafka instance on localhost:9092"
	go test -v -tags=integration ./...

# Run all tests
test: test-unit

# Build the package
build:
	@echo "Building messaging package..."
	go build ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Check for issues
vet:
	@echo "Running go vet..."
	go vet ./...

# Run all checks
check: fmt vet lint test

# Help
help:
	@echo "Available targets:"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests (requires Kafka)"
	@echo "  test          - Run all tests"
	@echo "  build         - Build the package"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Run linter"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  benchmark     - Run benchmarks"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  check         - Run all checks (fmt, vet, lint, test)"
	@echo "  help          - Show this help message"
