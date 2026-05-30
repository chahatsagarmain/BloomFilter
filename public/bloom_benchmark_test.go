package public_test

import (
	"fmt"
	"testing"

	"github.com/chahatsagarmain/bloomfilter/public"
)

func BenchmarkBitFilterInsert(b *testing.B) {
	filter, _ := public.BloomFactory("bit", 100000, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Insert(fmt.Sprintf("item-%d", i))
	}
}

func BenchmarkBitFilterContains(b *testing.B) {
	filter, _ := public.BloomFactory("bit", 100000, 4)
	filter.Insert("test-item")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Contains("test-item")
	}
}

func BenchmarkCountingFilterInsert(b *testing.B) {
	filter, _ := public.BloomFactory("counting", 100000, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Insert(fmt.Sprintf("item-%d", i))
	}
}

func BenchmarkCountingFilterContains(b *testing.B) {
	filter, _ := public.BloomFactory("counting", 100000, 4)
	filter.Insert("test-item")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Contains("test-item")
	}
}

func BenchmarkParallelBitFilterInsert(b *testing.B) {
	filter, _ := public.BloomFactory("bit", 100000, 4)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			filter.Insert(fmt.Sprintf("item-%d", i))
			i++
		}
	})
}

func BenchmarkParallelBitFilterContains(b *testing.B) {
	filter, _ := public.BloomFactory("bit", 100000, 4)
	filter.Insert("test-item")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			filter.Contains("test-item")
		}
	})
}

func BenchmarkParallelCountingFilterInsert(b *testing.B) {
	filter, _ := public.BloomFactory("counting", 100000, 4)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			filter.Insert(fmt.Sprintf("item-%d", i))
			i++
		}
	})
}

func BenchmarkParallelCountingFilterContains(b *testing.B) {
	filter, _ := public.BloomFactory("counting", 100000, 4)
	filter.Insert("test-item")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			filter.Contains("test-item")
		}
	})
}
