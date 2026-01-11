package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	packs "github.com/jagodawojcik/pack-calculator/internal/calculatepacks"
)


func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func packsHandler(w http.ResponseWriter, r *http.Request) {

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

		packSizes := []int{250, 500, 1000, 2000, 5000}
		result := packs.CalculatePacks(quantity, packSizes)

		response := map[string]any{
			"packs": result,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}


func main() {

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/packs", packsHandler)

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
