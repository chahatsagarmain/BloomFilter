package main

import (
	"fmt"
	"hash/fnv"
)

type bitbuckets struct {
	bitarray []bool 
	size     uint32
}

// NewBloomFilter creates a new filter with a specific size
func NewBloomFilter(size uint32) *bitbuckets {
	return &bitbuckets{
		bitarray: make([]bool, size),
		size:     size,
	}
}

func (b *bitbuckets) getHashIndices(s string) []uint32 {
	indices := []uint32{}
	// Let's use 3 hash functions (k=3)
	for i := 0; i < 3; i++ {
		h := fnv.New32a()
		// lets add a simple salt string + i
		h.Write([]byte(fmt.Sprintf("%s%d", s, i)))
		val := h.Sum32() % b.size
		indices = append(indices, val)
	}
	return indices
}

func (b *bitbuckets) insertString(s string) {
	indices := b.getHashIndices(s)
	for _, idx := range indices {
		b.bitarray[idx] = true
	}
	fmt.Printf("Inserted string: %s\n", s)
}

func (b *bitbuckets) checkString(s string) bool {
	indices := b.getHashIndices(s)
	for _, idx := range indices {
		if !b.bitarray[idx] {
			return false
		}
	}
	// If ALL bits were 1, it might be present
	return true
}

func main() {
	buckets := NewBloomFilter(1000)

	for {
		var mode int
		var input string
		fmt.Println("0: Insert string")
		fmt.Println("1: Check string")
		fmt.Print("Enter mode: ")
		fmt.Scanln(&mode)

		if mode == 0 {
			fmt.Print("Enter string to insert: ")
			fmt.Scanln(&input)
			buckets.insertString(input)
		} else if mode == 1 {
			fmt.Print("Enter string to check: ")
			fmt.Scanln(&input)
			result := buckets.checkString(input)
			if result {
				fmt.Printf("String '%s' is PROBABLY present (might be a false positive)\n", input)
			} else {
				fmt.Printf("String '%s' is DEFINITELY NOT present\n", input)
			}
		} else {
			fmt.Println("Wrong mode input")
		}
	}
}