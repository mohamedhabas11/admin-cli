# Define variables
PROJECT_NAME := admin-cli
GO := go

.PHONY: all test lint build clean

# Default target
all: test

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test ./internal/... -v

# Lint the code
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Build the project
build:
	@echo "Building $(PROJECT_NAME)..."
	@$(GO) build -o bin/$(PROJECT_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin

