package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	ctx1, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

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

	server := &http.Server{Addr: ":3000", Handler: r}

	serverErr := make(chan error, 1)

	go func() {

		err = server.ListenAndServe()

		if err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-ctx1.Done():
		log.Println("Shutdown signal received")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second)

		err = server.Shutdown(ctx)
		cancel()

		if err != nil {
			log.Println("Graceful shutdown failed: ", err)
		} else {
			log.Println("Server shutdown completed")
		}
	case er := <-serverErr:
		if er != nil && !errors.Is(er, http.ErrServerClosed) {
			log.Println("Server failed : ", er)
		}
	}

}
