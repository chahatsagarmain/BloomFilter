package bloom

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	f := New(1000, 3)

	// Ensure elements not inserted are not found (or at least, highly unlikely to collide)
	nonInserted := []string{"hello", "world", "golang", "programming"}
	for _, s := range nonInserted {
		if f.Contains(s) {
			t.Errorf("Expected Contains(%q) to be false, but was true", s)
		}
	}

	// Insert elements
	inserted := []string{"apple", "banana", "cherry"}
	for _, s := range inserted {
		f.Insert(s)
	}

	// Ensure inserted elements are always found (no false negatives)
	for _, s := range inserted {
		if !f.Contains(s) {
			t.Errorf("Expected Contains(%q) to be true, but was false", s)
		}
	}
}

func TestFillRatio(t *testing.T) {
	f := New(10, 1) // small size, 1 hash function to make math simple

	if f.FillRatio() != 0.0 {
		t.Errorf("Expected initial FillRatio to be 0.0, got %f", f.FillRatio())
	}

	// Insert a string that maps to a single index
	f.Insert("test")
	ratio := f.FillRatio()
	if ratio != 0.1 {
		t.Errorf("Expected FillRatio to be 0.1, got %f", ratio)
	}
}

func TestFalsePositiveRate(t *testing.T) {
	// Size: 500, k: 3
	// Insert 50 elements
	// Check false positive rate on 1000 non-inserted elements
	f := New(500, 3)

	inserted := make(map[string]bool)
	for i := 0; i < 50; i++ {
		s := fmt.Sprintf("item-%d", i)
		f.Insert(s)
		inserted[s] = true
	}

	// Check if all inserted elements are found
	for s := range inserted {
		if !f.Contains(s) {
			t.Fatalf("False negative detected: item %q was inserted but not found", s)
		}
	}

	// Check false positives
	falsePositives := 0
	numChecks := 1000
	for i := 0; i < numChecks; i++ {
		s := fmt.Sprintf("check-%d", i)
		if f.Contains(s) {
			falsePositives++
		}
	}

	fpRate := float64(falsePositives) / float64(numChecks)
	t.Logf("Empirical False Positive Rate: %f (%d/%d)", fpRate, falsePositives, numChecks)

	// Theoretical rate formula: (1 - e^(-k*n/m))^k
	// here m=500, n=50, k=3 => e^(-3 * 50 / 500) = e^(-0.3) approx 0.7408
	// (1 - 0.7408)^3 approx 0.2592^3 approx 0.0174 (1.7%)
	// We'll allow up to 5% to account for statistical variance.
	if fpRate > 0.05 {
		t.Errorf("False positive rate was too high: %f, expected < 0.05", fpRate)
	}
}

func TestConcurrency(t *testing.T) {
	f := New(10000, 4)

	const numGoroutines = 100
	const numInserts = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for j := 0; j < numInserts; j++ {
				f.Insert(fmt.Sprintf("item-%d-%d", gID, j))
			}
		}(i)
	}

	wg.Wait()

	// Verify that all inserted elements are present (no false negatives due to race conditions)
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numInserts; j++ {
			item := fmt.Sprintf("item-%d-%d", i, j)
			if !f.Contains(item) {
				t.Errorf("Expected Contains(%q) to be true, but got false", item)
			}
		}
	}
}
