.PHONY: build run test clean help

# Default target
all: test build

# Compile the interactive CLI tool binary
build:
	@echo "Building bloom-cli binary..."
	go build -o bin/bloom-cli cmd/bloom-cli/main.go

# Run the interactive CLI tool
run:
	go run cmd/bloom-cli/main.go

# Run all unit and statistical tests
test:
	@echo "Running unit tests..."
	go test -v ./...

# Remove compile artifacts and binaries
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin

# Show list of available commands
help:
	@echo "Available commands:"
	@echo "  make build  - Compile the bloom-cli binary"
	@echo "  make run    - Start the interactive CLI application"
	@echo "  make test   - Execute all package unit tests"
	@echo "  make clean  - Delete compiled binaries"
	@echo "  make help   - Display this help message"
