package app

import (
	"context"
	"errors"
	"fmt"
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

	log.Printf("[APP] Notification Service starting on port %d...", cfg.Server.Port)

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

	messageHandler := func(ctx context.Context, topic string, message []byte) error {
		log.Printf("[APP] Topic: %s, Message: %s", topic, string(message))
		if handler, ok := eventHandlers[topic]; ok {
			if err := handler.Handle(message); err != nil {
				log.Printf("[APP] Handler error for topic %s: %v", topic, err)
				return fmt.Errorf("handler error for topic %s: %w", topic, err)
			}
			return nil
		}
		log.Printf("[APP] No handler found for topic: %s", topic)
		return fmt.Errorf("no handler found for topic %s", topic)
	}
	
	consumer := infrastructure.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		topics,
		"notification-service",
		messageHandler,
	)

	consumer.Start(ctx)

	log.Println("[APP] Notification Service is running. Waiting for events...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("[APP] Notification Service shutting down...")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
