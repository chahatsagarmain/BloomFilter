package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/chahatsagarmain/bloomfilter/public"
)

type FPRRecord struct {
	InsertedElements int     `json:"inserted_elements"`
	FillRatio        float64 `json:"fill_ratio"`
	EmpiricalFPR     float64 `json:"empirical_fpr"`
	TheoreticalFPR   float64 `json:"theoretical_fpr"`
}

type KRecord struct {
	K              int     `json:"k"`
	EmpiricalFPR   float64 `json:"empirical_fpr"`
	TheoreticalFPR float64 `json:"theoretical_fpr"`
}

type BenchmarkResults struct {
	FPRVsElements []FPRRecord `json:"fpr_vs_elements"`
	FPRVsK        []KRecord   `json:"fpr_vs_k"`
}

func main() {
	outputFile := flag.String("output", "benchmark_results.json", "path to save benchmark results JSON")
	m := flag.Int("m", 100000, "Bloom filter size (bits)")
	kDefault := flag.Int("k", 4, "Number of hash functions for element filled test")
	nForKTest := flag.Int("n-ktest", 20000, "Number of elements for k test")
	numChecks := flag.Int("checks", 10000, "Number of false positive checks to perform")
	flag.Parse()

	fmt.Println("Starting Bloom Filter Benchmarks...")
	fmt.Printf("Configurations:\n  • Bit array size (m): %d\n  • k (default): %d\n  • n (for k-test): %d\n  • Checks: %d\n\n", *m, *kDefault, *nForKTest, *numChecks)

	// Run Experiment 1: FPR vs Elements filled
	fmt.Println("Running Experiment 1: FPR vs. Number of Elements filled...")
	fprVsElements := runFPRVsElements(*m, *kDefault, *numChecks)

	// Run Experiment 2: FPR vs k
	fmt.Println("Running Experiment 2: FPR vs. Number of Hash Functions (k)...")
	fprVsK := runFPRVsK(*m, *nForKTest, *numChecks)

	results := BenchmarkResults{
		FPRVsElements: fprVsElements,
		FPRVsK:        fprVsK,
	}

	// Write results to JSON
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", *outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("\nBenchmarks completed successfully. Results saved to: %s\n", *outputFile)
}

func runFPRVsElements(m, k, numChecks int) []FPRRecord {
	var records []FPRRecord

	// Target fill ratios from 0% to 98% in steps of 2%
	for pct := 0; pct <= 98; pct += 2 {
		targetRatio := float64(pct) / 100.0
		var n int
		if targetRatio == 0 {
			n = 0
		} else {
			// n = - (m / k) * ln(1 - targetRatio)
			nVal := -(float64(m) / float64(k)) * math.Log(1.0-targetRatio)
			n = int(math.Round(nVal))
		}

		// Create standard bit bloom filter
		filter, err := public.BloomFactory("bit", m, k)
		if err != nil {
			fmt.Printf("Error creating filter: %v\n", err)
			continue
		}

		// Insert n elements
		for i := 0; i < n; i++ {
			filter.Insert(fmt.Sprintf("insert-%d", i))
		}

		// Read fill ratio
		fillRatio := filter.FillRatio()

		// Measure empirical False Positive Rate (FPR)
		falsePositives := 0
		for i := 0; i < numChecks; i++ {
			// Query string that was never inserted
			checkStr := fmt.Sprintf("check-%d", i)
			if filter.Contains(checkStr) {
				falsePositives++
			}
		}
		empiricalFPR := float64(falsePositives) / float64(numChecks)

		// Calculate theoretical FPR
		// Formula: (1 - e^(-k * n / m))^k
		knOverM := float64(k) * float64(n) / float64(m)
		theoreticalFPR := math.Pow(1.0-math.Exp(-knOverM), float64(k))

		records = append(records, FPRRecord{
			InsertedElements: n,
			FillRatio:        fillRatio,
			EmpiricalFPR:     empiricalFPR,
			TheoreticalFPR:   theoreticalFPR,
		})

		fmt.Printf("  • Target Fill: %3d%% | Elements (n): %6d | Actual Fill: %6.2f%% | Empirical FPR: %6.4f | Theoretical FPR: %6.4f\n",
			pct, n, fillRatio*100, empiricalFPR, theoreticalFPR)
	}
	return records
}

func runFPRVsK(m, n, numChecks int) []KRecord {
	var records []KRecord

	for k := 1; k <= 16; k++ {
		// Create standard bit bloom filter with k hash functions
		filter, err := public.BloomFactory("bit", m, k)
		if err != nil {
			fmt.Printf("Error creating filter: %v\n", err)
			continue
		}

		// Insert n elements
		for i := 0; i < n; i++ {
			filter.Insert(fmt.Sprintf("insert-%d", i))
		}

		// Measure empirical False Positive Rate (FPR)
		falsePositives := 0
		for i := 0; i < numChecks; i++ {
			checkStr := fmt.Sprintf("check-%d", i)
			if filter.Contains(checkStr) {
				falsePositives++
			}
		}
		empiricalFPR := float64(falsePositives) / float64(numChecks)

		// Calculate theoretical FPR
		knOverM := float64(k) * float64(n) / float64(m)
		theoreticalFPR := math.Pow(1.0-math.Exp(-knOverM), float64(k))

		records = append(records, KRecord{
			K:              k,
			EmpiricalFPR:   empiricalFPR,
			TheoreticalFPR: theoreticalFPR,
		})

		fmt.Printf("  • Hash Functions (k): %2d | Empirical FPR: %6.4f | Theoretical FPR: %6.4f\n",
			k, empiricalFPR, theoreticalFPR)
	}
	return records
}
