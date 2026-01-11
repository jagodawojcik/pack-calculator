package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	packs "github.com/jagodawojcik/pack-calculator/internal/calculatepacks"
)

// Parses PACK_SIZES from environment variable
func readPackSizesFromEnv() []int {
	env := os.Getenv("PACK_SIZES")
	if strings.TrimSpace(env) == "" {
		panic("PACK_SIZES environment variable must be set and not empty")
	}

	parts := strings.Split(env, ",")
	packSizes := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		size, err := strconv.Atoi(p)
		if err != nil || size <= 0 {
			panic(fmt.Sprintf("invalid pack size '%s' in PACK_SIZES", p))
		}
		packSizes = append(packSizes, size)
	}

	return packSizes
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func packsHandler(packSizes []int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		origin := r.Header.Get("Origin")
		if origin == "https://jagodawojcik.github.io" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		quantityStr := r.URL.Query().Get("quantity")
		if quantityStr == "" {
			http.Error(w, "Quantity not specified", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(quantityStr)
		if err != nil || quantity <= 0 || quantity > 10000000 {
			http.Error(w, "Provide quantity between 1 and 10 000 000", http.StatusBadRequest)
			return
		}

		result := packs.CalculatePacks(quantity, packSizes)

		response := map[string]any{
			"packs": result,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	packSizes := readPackSizesFromEnv()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/packs", packsHandler(packSizes))

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
