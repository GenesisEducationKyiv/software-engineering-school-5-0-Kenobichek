package main

import (
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/external/sendgrid_email_api"
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/httpclient"
	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/scheduler"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"Weather-Forecast-API/internal/weather_provider"
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"net/http"
	"time"

	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/db"
)

func main() {
	//TODO: Move core logic to internal/app/app.go,
	// add config.go for env/settings, and invoke app.Run(ctx) from main.

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

	sgClient := sendgrid.NewSendClient(cfg.SendGrid.APIKey)
	sgNotifier := sendgrid_email_api.NewSendgridNotifier(sgClient, cfg)
	sgEmNotifier := notifier.NewSendGridEmailNotifier(&sgNotifier)

	httpClient := httpclient.New()

	geoSvc := openweather.NewOpenWeatherGeocodingService(cfg, httpClient)
	owAPI := openweather.NewOpenWeatherAPI(cfg, httpClient)

	weatherProvider := weather_provider.NewOpenWeatherProvider(&geoSvc, &owAPI)
	subscriptionService := subscription.NewSubscriptionService()
	notificationService := notification.NewNotificationService(&sgEmNotifier)

	newScheduler := scheduler.NewScheduler(cfg, &notificationService, &weatherProvider)
	go func() {
		_, err := newScheduler.Start()
		if err != nil {
			log.Printf("Error starting newScheduler: %v", err)
		}
	}()

	httpRouter := routes.NewHTTPRouter()

	weatherHandler := weather.NewWeatherHandler(&weatherProvider, 5*time.Second)
	subscribeHandler := subscribe.NewSubscribeHandler(subscriptionService, &notificationService)

	router := routes.NewRouter(&weatherHandler, &subscribeHandler, httpRouter)
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
