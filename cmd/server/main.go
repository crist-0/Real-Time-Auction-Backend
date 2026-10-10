package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Response struct {
	Status string `json:"status"`
}

func main() {

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		log.Fatal("Database connection failed : ", err)
	}

	err = pool.Ping(ctx)

	if err != nil {
		pool.Close()
		log.Fatal("Database connection refused : ", err)
	} else {
		log.Println("Database connection successfull : ")
	}

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status: "Ok",
		})
	})

	err = http.ListenAndServe(":3000", r)

	if err != nil {
		log.Fatal("server failed : ", err)
	}
}
