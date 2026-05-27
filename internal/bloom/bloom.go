package bloom

import (
	"hash/fnv"
)

// Filter represents a Bloom Filter data structure.
// It uses a slice of booleans for the bit representation.
type Filter struct {
	bitarray []bool
	size     uint32
	k        uint32 // number of hash functions
}

// New creates a new Bloom Filter with a specified size and k hash functions.
func New(size uint32, k uint32) *Filter {
	if k == 0 {
		k = 3 // default to 3 hash functions if not specified
	}
	return &Filter{
		bitarray: make([]bool, size),
		size:     size,
		k:        k,
	}
}

// getHashIndices returns the k index locations in the bit array for a given string
// using the Kirsch-Mitzenmacher optimization: hash_i = (hash1 + i * hash2) % size.
func (f *Filter) getHashIndices(s string) []uint32 {
	indices := make([]uint32, f.k)
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	hash64 := h.Sum64()

	// Split 64-bit hash into two 32-bit values
	hash1 := uint32(hash64 >> 32)
	hash2 := uint32(hash64)

	for i := uint32(0); i < f.k; i++ {
		indices[i] = (hash1 + i*hash2) % f.size
	}
	return indices
}

// Insert adds a string to the Bloom Filter by setting the corresponding bits to true.
func (f *Filter) Insert(s string) {
	indices := f.getHashIndices(s)
	for _, idx := range indices {
		f.bitarray[idx] = true
	}
}

// Contains checks if a string is in the Bloom Filter.
// Returns true if the string is probably present, and false if it is definitely not.
func (f *Filter) Contains(s string) bool {
	indices := f.getHashIndices(s)
	for _, idx := range indices {
		if !f.bitarray[idx] {
			return false
		}
	}
	return true
}

// FillRatio returns the percentage of bits in the bitarray set to true.
// This indicates the saturation of the filter.
func (f *Filter) FillRatio() float64 {
	if f.size == 0 {
		return 0.0
	}
	var setBits int
	for _, val := range f.bitarray {
		if val {
			setBits++
		}
	}
	return float64(setBits) / float64(f.size)
}

// Size returns the size of the Bloom Filter bit array.
func (f *Filter) Size() int {
	return int(f.size)
}

// K returns the number of hash functions used.
func (f *Filter) K() int {
	return int(f.k)
}
