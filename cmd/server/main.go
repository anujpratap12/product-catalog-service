package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"product-catalog-service/internal/handlers"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Server is running successfully",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", healthHandler)

	mux.HandleFunc("/request", handlers.RequestHandler)
	mux.HandleFunc("/stats", handlers.StatsHandler)
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			handlers.CreateProductHandler(w, r)
			return
		}

		if r.Method == http.MethodGet {
			handlers.ListProductsHandler(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {

		if strings.HasSuffix(r.URL.Path, "/media") {
			handlers.AddMediaHandler(w, r)
			return
		}

		handlers.GetProductHandler(w, r)
	})

	log.Println("Server running on port 8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
