package packs

import "math"

// Best solution for a given quantity
type solution struct {
	totalItems int
	numPacks   int
	lastPack   int
}
// CalculatePacks computes the optimal distribution of packs for a given order quantity.
func CalculatePacks(orderQuantity int, packSizes []int) map[int]int {
	if orderQuantity <= 0 {
		return map[int]int{}
	}

	// Exact quantity match shortcut
	for _, size := range packSizes {
		if orderQuantity == size {
			return map[int]int{size: 1}
		}
	}

	// DP table: best[target] = solution
	best := make(map[int]solution)
	best[0] = solution{totalItems: 0, numPacks: 0, lastPack: 0}

	for target := 1; target <= orderQuantity; target++ {
		bestItems := math.MaxInt
		bestPacks := math.MaxInt
		lastPack := 0

		for _, packSize := range packSizes {
			remaining := target - packSize
			if remaining < 0 {
				remaining = 0
			}
			
			// Check if there's a solution for the remaining quantity
			prev, ok := best[remaining]
			if !ok {
				continue
			}

			newItems := prev.totalItems + packSize
			newPacks := prev.numPacks + 1

			// First check for fewer excess items, then for fewer packs
			if newItems < bestItems || (newItems == bestItems && newPacks < bestPacks) {
				bestItems = newItems
				bestPacks = newPacks
				lastPack = packSize
			}
		}
		// Found a valid pack combination, save in DP table
		if lastPack > 0 {
			best[target] = solution{
				totalItems: bestItems,
				numPacks:   bestPacks,
				lastPack:   lastPack,
			}
		}
	}

	// Reconstruct distribution for the target order quantity
	result := make(map[int]int)
	current := orderQuantity
	for current > 0 {
		sol, ok := best[current]
		if !ok || sol.lastPack == 0 {
			break
		}
		result[sol.lastPack]++
		current -= sol.lastPack
	}

	return result
}
