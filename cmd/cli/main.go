package main

import (
	"fmt"
	"os"
	"strconv"

	packs "github.com/jagodawojcik/pack-calculator/internal/calculatepacks"
)

func main() {
	packSizes := []int{250, 500, 1000, 2000, 5000}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <quantity>")
		return
	}

	// Parse the quantity from CLI arg
	quantityStr := os.Args[1]
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		fmt.Println("Error: quantity must be a positive integer")
		return
	}

	// Calculate packs and print result
	result := packs.CalculatePacks(quantity, packSizes)

	fmt.Printf("Order quantity: %d\n", quantity)
	fmt.Println("Pack distribution:")
	for size, count := range result {
		fmt.Printf("  %d × %d\n", size, count)
	}
}
