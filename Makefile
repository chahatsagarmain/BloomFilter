.PHONY: build run test fmt clean benchmark help

# Default target
all: fmt test build

# Compile all binaries (interactive CLI tool and benchmark runner)
build:
	@echo "Building bloom-cli binary..."
	go build -o bin/bloom-cli cmd/bloom-cli/main.go
	@echo "Building bloom-benchmark binary..."
	go build -o bin/bloom-benchmark cmd/bloom-benchmark/main.go

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

# Run the benchmark experiment and plot results
benchmark:
	@echo "Running false positive rate benchmarks..."
	go run cmd/bloom-benchmark/main.go --output benchmark_results.json
	@echo "Plotting benchmark results using Matplotlib..."
	python benchmark/plot.py

# Remove compile artifacts, binaries, and benchmark results
clean:
	@echo "Cleaning build and benchmark artifacts..."
	rm -rf bin benchmark_results.json benchmark/*.png

# Show list of available commands
help:
	@echo "Available commands:"
	@echo "  make build      - Compile all binaries"
	@echo "  make run        - Start the interactive CLI application"
	@echo "  make test       - Execute all package unit tests"
	@echo "  make fmt        - Format Go source files using gofmt"
	@echo "  make benchmark  - Run benchmarks and generate false positive rate charts"
	@echo "  make clean      - Delete compiled binaries and generated benchmark charts"
	@echo "  make help       - Display this help message"

