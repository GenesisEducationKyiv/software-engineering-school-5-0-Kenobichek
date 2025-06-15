package main

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"Weather-Forecast-API/internal"
	"Weather-Forecast-API/internal/db"
	"Weather-Forecast-API/internal/scheduler"
	"github.com/go-chi/chi/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file.")
	}

	db.Init()
	db.RunMigrations(db.DataBase)

	go scheduler.StartScheduler()

	router := chi.NewRouter()
	internal.RegisterRoutes(router)

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server started at :%s\n", port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v\n", err)
	}
}
