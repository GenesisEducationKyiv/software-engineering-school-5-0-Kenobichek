package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/kafka"
	"api-gateway/internal/routes"
	"api-gateway/internal/weatherclient"

	"github.com/go-chi/chi/v5"
)

const (
	defaultRequestTimeout = 5 * time.Second
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	r := chi.NewRouter()

	publisher := kafka.NewPublisher(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("failed to close publisher: %v", err)
		}
	}()

	subscribeHandler := handlers.NewSubscribeHandler(publisher)

	weatherClient, err := weatherclient.New(cfg.WeatherServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to init weather client: %w", err)
	}

	weatherHandler := handlers.NewWeatherHandler(weatherClient)

	r.Route("/api", func(r chi.Router) {
		routes.RegisterRoutes(r, weatherHandler, subscribeHandler)
	})

	addr := ":" + cfg.Port

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: defaultRequestTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("api-gateway listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("api-gateway shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
