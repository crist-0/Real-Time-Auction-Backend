package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	Status string `json:"status"`
}

func main() {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status: "Ok",
		})
	})
	err := http.ListenAndServe(":3000", r)

	if err != nil {
		log.Fatal("server failed : ", err)
	}
}
