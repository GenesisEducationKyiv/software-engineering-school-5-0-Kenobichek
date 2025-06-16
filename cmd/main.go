package main

import (
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/scheduler"
	"errors"
	"log"
	"net/http"
	"time"

	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	database, err := db.Init(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err = db.RunMigrations(database); err != nil {
		log.Printf("Error running database migrations: %v", err)
		return
	}

	go scheduler.StartScheduler()

	router := chi.NewRouter()
	routes.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server started at %s\n", cfg.GetServerAddress())

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %v\n", err)
	}
}
