# Simple Bloom Filter in Go

A lightweight, clean, and modular implementation of a **Bloom Filter** written in Go. This project demonstrates how probabilistic data structures work using multiple hash functions (simulated via salted FNV-1a) to store and check string existence efficiently.

## Project Structure

This project follows idiomatic Go practices:

```
bloomfilter/
├── go.mod                     # Go module definition
├── internal/
│   └── bloom/
│       ├── bloom.go           # Internal Bloom Filter library logic
│       └── bloom_test.go      # Comprehensive unit & statistical tests
├── cmd/
│   └── bloom-cli/
│       └── main.go            # CLI interactive application
└── README.md                  # Project documentation
```

## Features

- **Encapsulated Library**: Core logic is package-isolated inside `internal/bloom` to keep it module-private.
- **Configurable Parameters**: Parameterized filter size and number of hash functions ($k$).
- **Interactive CLI**: Menu-driven interface with terminal coloring, custom configuration, and saturation metrics (Fill Ratio).
- **Rigorous Tests**: Tests covering core operations, verification of 0% false negative rates, and validating the empirical false positive rate against mathematical expectations.

## Prerequisites

- [Go](https://go.dev/dl/) (version 1.20 or later recommended) installed.

## How to Run the CLI

To start the interactive command-line application:

```bash
go run cmd/bloom-cli/main.go
```

### CLI Options
Upon running the CLI, you can choose to use default configurations (size 1000, 3 hash functions) or define your own. The application offers options to:
1. **Insert strings** (displays the updated fill ratio after insertion).
2. **Check strings** (reports whether a string is *definitely not present* or *probably present*).
3. **View saturation statistics** (warns if the saturation gets high, leading to elevated false positive rates).
4. **Exit** the application.

## Running Tests

To execute the unit and statistical verification tests, run:

```bash
go test -v ./...
```
