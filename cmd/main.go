package main

import (
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/scheduler"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"Weather-Forecast-API/internal/weather_provider"
	"errors"
	"log"
	"net/http"
	"time"

	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/db"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	config.Usage()

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

	weatherProvider := weather_provider.NewOpenWeatherProvider(cfg.OpenWeather.APIKey)
	subscriptionService := subscription.NewSubscriptionService()
	notificationService := notification.NewNotificationService(cfg)

	newScheduler := scheduler.NewScheduler(cfg, notificationService, &weatherProvider)
	go func() {
		_, err := newScheduler.Start()
		if err != nil {
			log.Printf("Error starting newScheduler: %v", err)
		}
	}()

	router := routes.NewRouter(subscriptionService, notificationService)
	router.RegisterRoutes()

	srv := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router.GetRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server started at %s\n", cfg.GetServerAddress())

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %v\n", err)
	}
}
