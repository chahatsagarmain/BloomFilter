.PHONY: build run test fmt clean help

# Default target
all: fmt test build

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

# Format all Go source files
fmt:
	@echo "Formatting Go source files..."
	gofmt -s -w .

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
	@echo "  make fmt    - Format Go source files using gofmt"
	@echo "  make clean  - Delete compiled binaries"
	@echo "  make help   - Display this help message"
