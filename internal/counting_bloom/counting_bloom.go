package countingbloom

import (
	"hash/fnv"
	"sync"
)

type CountingFilter struct {
	sync.RWMutex
	array []uint8
	size  uint32
	k     uint32
}

func NewFilter(size uint32, k uint32) *CountingFilter {
	if k == 0 {
		k = 3
	}
	return &CountingFilter{
		array: make([]uint8, size),
		size:  size,
		k:     k,
	}
}

func (f *CountingFilter) getHashIndices(s string) []uint32 {
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

func (f *CountingFilter) Insert(s string) {
	indices := f.getHashIndices(s)
	f.Lock()
	defer f.Unlock()
	for _, idx := range indices {
		if f.array[idx] < 255 {
			f.array[idx]++
		}
	}
}

func (f *CountingFilter) Contains(s string) bool {
	indices := f.getHashIndices(s)
	f.RLock()
	defer f.RUnlock()
	for _, idx := range indices {
		if f.array[idx] == 0 {
			return false
		}
	}
	return true
}

func (f *CountingFilter) FillRatio() float64 {
	f.RLock()
	defer f.RUnlock()

	if f.size == 0 {
		return 0.0
	}

	setBits := 0
	for _, val := range f.array {
		if val > 0 {
			setBits++
		}
	}

	return float64(setBits) / float64(f.size)
}

func (f *CountingFilter) Delete(s string) {
	indices := f.getHashIndices(s)
	f.Lock()
	defer f.Unlock()
	for _, idx := range indices {
		if f.array[idx] > 0 {
			f.array[idx]--
		}
	}
}

func (f *CountingFilter) Size() int {
	return int(f.size)
}

func (f *CountingFilter) K() int {
	return int(f.k)
}
