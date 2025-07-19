package app

import (
	"context"
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

	"github.com/go-chi/chi/v5"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	r := chi.NewRouter()

	publisher := kafka.NewPublisher(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer publisher.Close()

	subscribeHandler := handlers.NewSubscribeHandler(publisher)

	weatherHandler, err := handlers.NewWeatherHandler(cfg.WeatherServiceAddr)
	if err != nil {
		log.Fatalf("failed to init WeatherHandler: %v", err)
	}
	
	r.Route("/api", func(r chi.Router) {
		routes.RegisterRoutes(r, weatherHandler, subscribeHandler)
	})

	addr := ":" + cfg.Port

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
