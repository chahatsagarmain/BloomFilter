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

## How It Works (Step-by-Step)

A Bloom Filter is a space-efficient probabilistic data structure. Here is how the underlying mechanics operate step-by-step:

### 1. Initialization
When you create a Bloom Filter (using `public.BloomFactory`), you specify:
* **Size ($m$)**: The number of slots (bits in a standard filter, or 8-bit counters in a counting filter).
* **Hash Functions ($k$)**: The number of hash indices generated for each string.

### 2. Insertion Step-by-Step
To insert a string (e.g., `"golang"`):
1. **Generate Hash Values**: The string is hashed using the FNV-1a hashing algorithm to produce a 64-bit hash value.
2. **Split Hash**: The 64-bit hash is split into two 32-bit values ($hash_1$ and $hash_2$).
3. **Kirsch-Mitzenmacher Optimization**: The filter generates $k$ unique slot indices using the formula:
   $$index_i = (hash_1 + i \times hash_2) \pmod{size}$$
   This allows us to simulate $k$ independent hash functions in $O(1)$ time with only a single hashing pass.
4. **Update Slots**:
   * **Standard Filter**: Sets the bit at each generated index to `true`.
   * **Counting Filter**: Increments the 8-bit counter at each generated index. The counters use **saturating arithmetic** (stopping at `255` instead of wrapping around to `0` to prevent overflow).

### 3. Membership Lookup Step-by-Step
To check if a string exists (e.g., `"golang"`):
1. **Generate Hash Indices**: Calculate the same $k$ indices for the string using steps 1 to 3 above.
2. **Scan Slots**:
   * **Standard Filter**: Check if the bit at every index is `true`.
   * **Counting Filter**: Check if the counter at every index is greater than `0`.
3. **Evaluation**:
   * If **any** index contains `false` (standard) or `0` (counting), the item is **DEFINITELY NOT** in the set (0% false negative rate).
   * If **all** indices contain `true` or $> 0$, the item is **PROBABLY** in the set. There is a small chance of a **false positive** if other elements' hashes have overlapped to set all these slots.

### 4. Deletion Step-by-Step (Counting Bloom Filter only)
Standard Bloom Filters do not support deletion because clearing bits could accidentally delete other overlapping elements. Counting Bloom Filters solve this:
1. **Generate Hash Indices**: Calculate the $k$ indices for the string to delete.
2. **Decrement Counters**: Decrement the counter at each of these indices (stopping at `0` to prevent underflow). Since counters are decremented rather than setting a bit to false, overlapping elements remain correctly represented if their counter is still $> 0$.

### 5. Update Step-by-Step (Counting Bloom Filter only)
To atomically replace an old element with a new element:
1. **Calculate Indices**: Generate $k$ indices for the old string, and another $k$ indices for the new string.
2. **Atomic Lock**: Acquire a write lock to ensure thread safety.
3. **Decrement Old**: Decrement counters at the indices of the old string.
4. **Increment New**: Increment counters at the indices of the new string.
5. **Release Lock**: Release the write lock.

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

### 2. Counting Bloom Filter (with Deletion & Update)

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

	// Update the element (only supported by counting Bloom filters)
	if deletable, ok := filter.(public.DeletableBloomFilter); ok {
		deletable.Update("golang", "rust")
	}

	fmt.Println("Contains golang after update:", filter.Contains("golang")) // false
	fmt.Println("Contains rust after update:", filter.Contains("rust"))     // true
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
2. **Counting Bloom Filter** (enables the delete and update options in the menu)

### Example CLI Sessions

#### 1. Standard Bit Bloom Filter Session

```text
=============================================
      Interactive Bloom Filter CLI
=============================================
Choose Bloom Filter type (1: Standard Bit, 2: Counting) [Default: 1]: 1
Do you want to configure filter size and hash functions? (y/N): n

✔ Bloom Filter initialized successfully!
  • Type: bit
  • Size: 1000 bits
  • Hash Functions (k): 3

---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Exit
Enter choice: 1
Enter string to insert: hello
✔ Successfully inserted string: "hello"
  Current Fill Ratio: 0.30%
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Exit
Enter choice: 1
Enter string to insert: world
✔ Successfully inserted string: "world"
  Current Fill Ratio: 0.60%
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Exit
Enter choice: 2
Enter string to check: hello
★ String "hello" is PROBABLY in the filter (might be a false positive).
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Exit
Enter choice: 2
Enter string to check: golang
✖ String "golang" is DEFINITELY NOT in the filter.
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Exit
Enter choice: 3
Filter Saturation (Fill Ratio): 0.6000%
✔ Good saturation level.
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Exit
Enter choice: 4

Goodbye!
```

#### 2. Counting Bloom Filter Session (with deletion and update support)

```text
=============================================
      Interactive Bloom Filter CLI
=============================================
Choose Bloom Filter type (1: Standard Bit, 2: Counting) [Default: 1]: 2
Do you want to configure filter size and hash functions? (y/N): n

✔ Bloom Filter initialized successfully!
  • Type: counting
  • Size: 1000 bits
  • Hash Functions (k): 3

---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 1
Enter string to insert: golang
✔ Successfully inserted string: "golang"
  Current Fill Ratio: 0.30%
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 2
Enter string to check: golang
★ String "golang" is PROBABLY in the filter (might be a false positive).
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 5
Enter old string to replace: golang
Enter new string: rust
✔ Successfully updated string from "golang" to "rust"
  Current Fill Ratio: 0.30%
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 2
Enter string to check: golang
✖ String "golang" is DEFINITELY NOT in the filter.
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 2
Enter string to check: rust
★ String "rust" is PROBABLY in the filter (might be a false positive).
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 3
Filter Saturation (Fill Ratio): 0.3000%
✔ Good saturation level.
---------------------------------------------
Select an option:
  1) Insert string
  2) Check string
  3) View saturation (Fill Ratio)
  4) Delete string
  5) Update string
  6) Exit
Enter choice: 6

Goodbye!
```

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
