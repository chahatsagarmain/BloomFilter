package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chahatsagarmain/bloomfilter/public"
)

// ANSI color codes for enhanced terminal UX
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(Bold + Cyan + "=============================================" + Reset)
	fmt.Println(Bold + Cyan + "      Interactive Bloom Filter CLI" + Reset)
	fmt.Println(Bold + Cyan + "=============================================" + Reset)

	bloomType := "bit"
	fmt.Printf("Choose Bloom Filter type (1: Standard Bit, 2: Counting) [Default: 1]: ")
	typeChoice, _ := reader.ReadString('\n')
	typeChoice = strings.TrimSpace(typeChoice)
	if typeChoice == "2" {
		bloomType = "counting"
	}

	// Optional custom configuration on startup
	size := 1000
	k := 3

	fmt.Printf("Do you want to configure filter size and hash functions? (y/N): ")
	configChoice, _ := reader.ReadString('\n')
	configChoice = strings.TrimSpace(strings.ToLower(configChoice))

	if configChoice == "y" || configChoice == "yes" {
		fmt.Print("Enter filter size (default 1000): ")
		sizeStr, _ := reader.ReadString('\n')
		sizeStr = strings.TrimSpace(sizeStr)
		if val, err := strconv.Atoi(sizeStr); err == nil && val > 0 {
			size = val
		}

		fmt.Print("Enter number of hash functions (default 3): ")
		kStr, _ := reader.ReadString('\n')
		kStr = strings.TrimSpace(kStr)
		if val, err := strconv.Atoi(kStr); err == nil && val > 0 {
			k = val
		}
	}

	filter, err := public.BloomFactory(bloomType, size, k)
	if err != nil {
		fmt.Printf(Red + "✗ Error initializing Bloom Filter: %s\n" + Reset, err.Error())
		return
	}

	fmt.Println(Bold + Green + "\n✔ Bloom Filter initialized successfully!" + Reset)
	fmt.Printf("  • Type: %s\n", bloomType)
	fmt.Printf("  • Size: %d bits\n", size)
	fmt.Printf("  • Hash Functions (k): %d\n\n", k)

	deletableFilter, isDeletable := filter.(public.DeletableBloomFilter)

	for {
		fmt.Println(Bold + Blue + "---------------------------------------------" + Reset)
		fmt.Println("Select an option:")
		fmt.Println(Cyan + "  1)" + Reset + " Insert string")
		fmt.Println(Cyan + "  2)" + Reset + " Check string")
		fmt.Println(Cyan + "  3)" + Reset + " View saturation (Fill Ratio)")
		if isDeletable {
			fmt.Println(Cyan + "  4)" + Reset + " Delete string")
			fmt.Println(Red + "  5)" + Reset + " Exit")
		} else {
			fmt.Println(Red + "  4)" + Reset + " Exit")
		}
		fmt.Print(Bold + "Enter choice: " + Reset)

		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Exiting...")
			break
		}
		choiceStr = strings.TrimSpace(choiceStr)

		switch choiceStr {
		case "1":
			fmt.Print("Enter string to insert: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "" {
				fmt.Println(Yellow + "⚠ Warning: Cannot insert an empty string." + Reset)
				continue
			}
			filter.Insert(input)
			fmt.Printf(Green+"✔ Successfully inserted string: "+Reset+"\"%s\"\n", input)
			fmt.Printf("  Current Fill Ratio: %.2f%%\n", filter.FillRatio()*100)

		case "2":
			fmt.Print("Enter string to check: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "" {
				fmt.Println(Yellow + "⚠ Warning: Cannot check an empty string." + Reset)
				continue
			}
			present := filter.Contains(input)
			if present {
				fmt.Printf(Yellow+"★ String \"%s\" is PROBABLY in the filter "+Reset+"(might be a false positive).\n", input)
			} else {
				fmt.Printf(Green+"✖ String \"%s\" is DEFINITELY NOT in the filter.\n"+Reset, input)
			}

		case "3":
			ratio := filter.FillRatio()
			fmt.Printf("Filter Saturation (Fill Ratio): "+Bold+Yellow+"%.4f%%\n"+Reset, ratio*100)
			// Give a hint if it's getting too full
			if ratio > 0.5 {
				fmt.Println(Red + "⚠ Warning: High saturation (>50%). False positive rate will increase." + Reset)
			} else {
				fmt.Println(Green + "✔ Good saturation level." + Reset)
			}

		case "4":
			if isDeletable {
				fmt.Print("Enter string to delete: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input == "" {
					fmt.Println(Yellow + "⚠ Warning: Cannot delete an empty string." + Reset)
					continue
				}
				deletableFilter.Delete(input)
				fmt.Printf(Green+"✔ Successfully deleted string: "+Reset+"\"%s\"\n", input)
				fmt.Printf("  Current Fill Ratio: %.2f%%\n", filter.FillRatio()*100)
			} else {
				fmt.Println(Bold + Green + "\nGoodbye!" + Reset)
				return
			}

		case "5":
			if isDeletable {
				fmt.Println(Bold + Green + "\nGoodbye!" + Reset)
				return
			} else {
				fmt.Println(Red + "✗ Invalid option. Please enter a valid number." + Reset)
			}

		case "exit", "quit":
			fmt.Println(Bold + Green + "\nGoodbye!" + Reset)
			return

		default:
			if isDeletable {
				fmt.Println(Red + "✗ Invalid option. Please enter 1, 2, 3, 4, or 5." + Reset)
			} else {
				fmt.Println(Red + "✗ Invalid option. Please enter 1, 2, 3, or 4." + Reset)
			}
		}
	}
}
