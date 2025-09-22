# Makefile for Go project

.PHONY: test test-coverage test-verbose clean build run help

# Default target
all: test build

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Generate detailed coverage report
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	go test -race ./...

# Run benchmarks
benchmark:
	go test -bench=. ./...

# Build the application
build:
	go build -o bin/main main.go

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f bin/main
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter (if golangci-lint is installed)
lint:
	golangci-lint run

# Run all checks (format, lint, test)
check: fmt lint test

# Help
help:
	@echo "Available targets:"
	@echo "  test              - Run tests"
	@echo "  test-verbose      - Run tests with verbose output"
	@echo "  test-coverage     - Run tests with coverage"
	@echo "  test-coverage-html- Generate HTML coverage report"
	@echo "  test-race         - Run tests with race detection"
	@echo "  benchmark         - Run benchmarks"
	@echo "  build             - Build the application"
	@echo "  run               - Run the application"
	@echo "  clean             - Clean build artifacts"
	@echo "  fmt               - Format code"
	@echo "  lint              - Run linter"
	@echo "  check             - Run all checks"
	@echo "  help              - Show this help"
