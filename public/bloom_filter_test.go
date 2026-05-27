package public

import (
	"testing"
)

func TestBloomFactoryInvalidParameters(t *testing.T) {
	_, err := BloomFactory("bit", 0, 3)
	if err == nil {
		t.Error("Expected error for size = 0, got nil")
	}

	_, err = BloomFactory("bit", 1000, -1)
	if err == nil {
		t.Error("Expected error for k = -1, got nil")
	}

	_, err = BloomFactory("invalid_type", 1000, 3)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

func TestBitBloomFilter(t *testing.T) {
	filter, err := BloomFactory("bit", 1000, 3)
	if err != nil {
		t.Fatalf("Failed to create bit bloom filter: %v", err)
	}

	// Verify metadata getters
	if filter.Size() != 1000 {
		t.Errorf("Expected Size() to be 1000, got %d", filter.Size())
	}
	if filter.K() != 3 {
		t.Errorf("Expected K() to be 3, got %d", filter.K())
	}

	// Verify basic operations
	if filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be false initially")
	}

	filter.Insert("hello")
	if !filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be true after Insert")
	}

	ratio := filter.FillRatio()
	if ratio <= 0.0 || ratio >= 1.0 {
		t.Errorf("Expected FillRatio to be between 0 and 1, got %f", ratio)
	}
}

func TestCountingBloomFilter(t *testing.T) {
	filter, err := BloomFactory("counting", 1000, 3)
	if err != nil {
		t.Fatalf("Failed to create counting bloom filter: %v", err)
	}

	// Verify metadata getters
	if filter.Size() != 1000 {
		t.Errorf("Expected Size() to be 1000, got %d", filter.Size())
	}
	if filter.K() != 3 {
		t.Errorf("Expected K() to be 3, got %d", filter.K())
	}

	// Verify basic operations
	if filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be false initially")
	}

	filter.Insert("hello")
	if !filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be true after Insert")
	}

	ratio := filter.FillRatio()
	if ratio <= 0.0 || ratio >= 1.0 {
		t.Errorf("Expected FillRatio to be between 0 and 1, got %f", ratio)
	}

	// Verify DeletableBloomFilter type assertion and deletion behavior
	deletable, ok := filter.(DeletableBloomFilter)
	if !ok {
		t.Fatal("Counting Bloom filter should implement DeletableBloomFilter interface")
	}

	deletable.Delete("hello")
	if filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be false after Delete")
	}
}

func TestCountingBloomFilterSaturatingArithmetic(t *testing.T) {
	// Small size, 1 hash function
	filter, err := BloomFactory("counting", 10, 1)
	if err != nil {
		t.Fatalf("Failed to create counting bloom filter: %v", err)
	}

	deletable, ok := filter.(DeletableBloomFilter)
	if !ok {
		t.Fatal("Counting Bloom filter should implement DeletableBloomFilter interface")
	}

	// Insert "hello" 300 times (exceeds max uint8 of 255)
	for i := 0; i < 300; i++ {
		filter.Insert("hello")
	}

	if !filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be true")
	}

	// Delete "hello" 300 times
	for i := 0; i < 300; i++ {
		deletable.Delete("hello")
	}

	// Due to saturating arithmetic (counter stops at 255 and stops at 0),
	// it should be successfully deleted and not wrap around or underflow
	if filter.Contains("hello") {
		t.Error("Expected Contains('hello') to be false after deletion matches insertion count")
	}
}

func TestCountingBloomFilterUpdate(t *testing.T) {
	filter, err := BloomFactory("counting", 1000, 3)
	if err != nil {
		t.Fatalf("Failed to create counting bloom filter: %v", err)
	}

	deletable, ok := filter.(DeletableBloomFilter)
	if !ok {
		t.Fatal("Counting Bloom filter should implement DeletableBloomFilter interface")
	}

	// Insert "golang"
	filter.Insert("golang")
	if !filter.Contains("golang") {
		t.Error("Expected Contains('golang') to be true")
	}
	if filter.Contains("rust") {
		t.Error("Expected Contains('rust') to be false")
	}

	// Update "golang" to "rust"
	deletable.Update("golang", "rust")

	// "golang" should be deleted (not present anymore)
	if filter.Contains("golang") {
		t.Error("Expected Contains('golang') to be false after update")
	}
	// "rust" should be inserted (present)
	if !filter.Contains("rust") {
		t.Error("Expected Contains('rust') to be true after update")
	}
}
