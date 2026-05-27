
package public

import (
	"fmt"

	"github.com/chahatsagarmain/bloomfilter/internal/bloom"
	countingbloom "github.com/chahatsagarmain/bloomfilter/internal/counting_bloom"
)

// BloomFilter defines the public behavior of a Bloom Filter.
type BloomFilter interface {
	Insert(s string)
	Contains(s string) bool
	FillRatio() float64
	Size() int
	K() int
}

// DeletableBloomFilter defines a Bloom Filter that supports element deletion.
type DeletableBloomFilter interface {
	BloomFilter
	Delete(s string)
}

// BloomFactory creates and returns a BloomFilter of the specified type ("bit" or "counting").
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
