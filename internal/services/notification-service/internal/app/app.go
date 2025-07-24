package app

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/config"
	"notification-service/internal/domain"
	"notification-service/internal/handlers"
	"notification-service/internal/infrastructure"
	"notification-service/internal/notifier"

	"github.com/sendgrid/sendgrid-go"
)

type eventHandlerManager interface {
	Handle(message []byte) error
}

func Run(ctx context.Context) error {
	cfg, err := config.MustLoad()
	if err != nil {
		return errors.New("failed to load config: " + err.Error())
	}

	log.Printf("Notification Service starting on port %d...", cfg.Server.Port)

	sgClient := sendgrid.NewSendClient(cfg.SendGrid.APIKey)
	sendgridNotifier := infrastructure.NewSendgridNotifier(sgClient, cfg.SendGrid.SenderName, cfg.SendGrid.SenderEmail)

	templateRepo := domain.NewTemplateRepository()

	notificationService := notifier.NewService(sendgridNotifier, templateRepo)

	eventHandlers := map[string]eventHandlerManager{
		"weather.updated":        handlers.NewWeatherUpdateHandler(notificationService),
		"subscription.confirmed": handlers.NewSubscriptionConfirmedHandler(notificationService),
		"subscription.cancelled": handlers.NewSubscriptionCancelledHandler(notificationService),
	}

	topics := make([]string, 0, len(eventHandlers))
	for topic := range eventHandlers {
		topics = append(topics, topic)
	}

	consumer := infrastructure.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		topics,
		"notification-service",
		func(topic string, message []byte) error {
			log.Printf("[HANDLER] Topic: %s, Message: %s", topic, string(message))
			if handler, ok := eventHandlers[topic]; ok {
				return handler.Handle(message)
			}
			log.Printf("No handler found for topic: %s", topic)
			return nil
		},
	)

	consumer.Start(ctx)

	log.Println("Notification Service is running. Waiting for events...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("Notification Service shutting down...")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
