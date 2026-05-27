// Package public provides the public interface and factory for constructing
// standard bit-array and counting Bloom Filters.
//
// Bloom filters are space-efficient probabilistic data structures used to test
// whether an element is a member of a set. False positive matches are possible,
// but false negatives are not.
package public

import (
	"fmt"

	"github.com/chahatsagarmain/bloomfilter/internal/bloom"
	countingbloom "github.com/chahatsagarmain/bloomfilter/internal/counting_bloom"
)

// BloomFilter defines the public behavior of a Bloom Filter.
type BloomFilter interface {
	// Insert adds a string to the Bloom Filter.
	Insert(s string)

	// Contains checks if a string is probably in the Bloom Filter (returns true),
	// or definitely not in the Bloom Filter (returns false).
	Contains(s string) bool

	// FillRatio returns the percentage of bits or buckets in the filter currently set.
	FillRatio() float64

	// Size returns the bit or bucket capacity of the Bloom Filter.
	Size() int

	// K returns the number of hash functions used.
	K() int
}

// DeletableBloomFilter defines a Bloom Filter that supports element deletion.
// Counting Bloom Filters implement this interface.
type DeletableBloomFilter interface {
	BloomFilter

	// Delete removes a string from the Bloom Filter.
	Delete(s string)
}

// BloomFactory creates and returns a BloomFilter of the specified type ("bit" or "counting").
// It returns an error if the size or k is less than or equal to zero, or if the bloomType is unsupported.
func BloomFactory(bloomType string, size, k int) (BloomFilter, error) {
	if size <= 0 || k <= 0 {
		return nil, fmt.Errorf("size and k must be greater than zero")
	}

	switch bloomType {
	case "bit":
		return bloom.New(uint32(size), uint32(k)), nil
	case "counting":
		return countingbloom.NewFilter(uint32(size), uint32(k)), nil
	default:
		return nil, fmt.Errorf("choose between bit / counting bloom filter")
	}
}
