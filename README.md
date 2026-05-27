# Simple Bloom Filter in Go

A lightweight, clean, and modular implementation of a **Bloom Filter** written in Go. This library provides both a standard bit-array Bloom Filter and a Counting Bloom Filter (which supports element deletion).

It follows Go software design best practices by exposing a clean, simple public API while keeping concrete implementations private within the `internal` package.

---

## Project Structure

This project follows idiomatic Go folder layout:

```
bloomfilter/
├── LICENSE                    # MIT License
├── Makefile                   # CLI and testing shortcuts
├── README.md                  # Project documentation
├── go.mod                     # Go module definition
├── cmd/
│   └── bloom-cli/
│       └── main.go            # Interactive terminal application
├── internal/
│   ├── bloom/
│   │   ├── bloom.go           # Standard bit-array Bloom Filter logic
│   │   └── bloom_test.go      # Bit Bloom Filter unit tests
│   └── counting_bloom/
│       └── counting_bloom.go  # Counting Bloom Filter logic
└── public/
    ├── bloom_filter.go        # Public interfaces (BloomFilter, DeletableBloomFilter) and factory
    ├── bloom_filter_test.go   # Package unit & statistical tests
    └── example_test.go        # Runnable GoDoc usage examples
```

---

## Features

- **Encapsulated Library**: Implementation details are private within `internal/`, exposing only the interfaces through the `public/` package.
- **Factory Pattern**: Easily construct standard bit-array filters or counting filters through `public.BloomFactory`.
- **Element Deletion**: Counting Bloom Filter allows deleting elements safely using 8-bit counters with saturating arithmetic (prevents overflow/underflow).
- **Interactive CLI**: Interactive terminal application to play with both filter types, insert/check strings, check saturation ratio, and delete items.
- **High Performance**: Salted FNV-1a hashing using the Kirsch-Mitzenmacher optimization ($hash_i = (hash_1 + i \times hash_2) \pmod{size}$) to achieve $O(k)$ operations.

---

## Installation

Add this package to your Go project:

```bash
go get github.com/chahatsagarmain/bloomfilter
```

---

## Quick Start

Here is a quick example of how to use the library in your own Go codebase.

### 1. Standard Bit-Array Bloom Filter

```go
package main

import (
	"fmt"
	"github.com/chahatsagarmain/bloomfilter/public"
)

func main() {
	// Initialize a standard bit bloom filter (size: 1000, hash functions k: 3)
	filter, err := public.BloomFactory("bit", 1000, 3)
	if err != nil {
		panic(err)
	}

	// Insert elements
	filter.Insert("hello")
	filter.Insert("world")

	// Check elements
	fmt.Println(filter.Contains("hello")) // Output: true
	fmt.Println(filter.Contains("go"))    // Output: false (definitely not in set)
}
```

### 2. Counting Bloom Filter (with Deletion)

```go
package main

import (
	"fmt"
	"github.com/chahatsagarmain/bloomfilter/public"
)

func main() {
	// Initialize a counting Bloom filter
	filter, err := public.BloomFactory("counting", 1000, 3)
	if err != nil {
		panic(err)
	}

	filter.Insert("golang")
	fmt.Println("Contains golang:", filter.Contains("golang")) // true

	// Check and delete if the filter supports deletion
	if deletable, ok := filter.(public.DeletableBloomFilter); ok {
		deletable.Delete("golang")
	}

	fmt.Println("Contains golang after delete:", filter.Contains("golang")) // false
}
```

---

## Interactive CLI

To run the interactive command-line application:

```bash
make run
```
or
```bash
go run cmd/bloom-cli/main.go
```

Upon launching, choose between:
1. **Standard Bit Bloom Filter**
2. **Counting Bloom Filter** (enables the delete option in the menu)

---

## Running Tests

To execute the unit, statistical verification, and runnable examples:

```bash
make test
```
or
```bash
go test -v ./...
```
