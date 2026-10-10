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
	"github.com/joho/godotenv"
)

type Response struct {
	Status string `json:"status"`
}

type App struct {
	DB *pgxpool.Pool
}

func (a *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status: "Ok",
	})
}

func main() {

	er := godotenv.Load()

	if er != nil {
		log.Fatal("Error loading .env file")
	}

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

	defer pool.Close()

	app := App{
		DB: pool,
	}

	r := chi.NewRouter()

	r.Get("/health", app.HealthHandler)

	err = http.ListenAndServe(":3000", r)

	if err != nil {
		log.Fatal("server failed : ", err)
	}
}
