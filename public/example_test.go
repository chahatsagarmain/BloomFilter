package public_test

import (
	"fmt"

	"github.com/chahatsagarmain/bloomfilter/public"
)

func ExampleBloomFactory_bit() {
	// Create a standard bit Bloom filter
	filter, err := public.BloomFactory("bit", 1000, 3)
	if err != nil {
		panic(err)
	}

	filter.Insert("hello")
	fmt.Println("Contains 'hello':", filter.Contains("hello"))
	fmt.Println("Contains 'world':", filter.Contains("world"))

	// Output:
	// Contains 'hello': true
	// Contains 'world': false
}

func ExampleBloomFactory_counting() {
	// Create a counting Bloom filter
	filter, err := public.BloomFactory("counting", 1000, 3)
	if err != nil {
		panic(err)
	}

	filter.Insert("golang")
	fmt.Println("Contains 'golang':", filter.Contains("golang"))

	// Delete the item (only supported by counting Bloom filters)
	if deletable, ok := filter.(public.DeletableBloomFilter); ok {
		deletable.Delete("golang")
	}
	fmt.Println("Contains 'golang' after delete:", filter.Contains("golang"))

	// Output:
	// Contains 'golang': true
	// Contains 'golang' after delete: false
}
