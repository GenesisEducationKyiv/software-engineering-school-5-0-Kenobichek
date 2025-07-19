package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/config"
	"notification-service/domain"
	"notification-service/infrastructure"
	"notification-service/notifier"

	"github.com/sendgrid/sendgrid-go"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Notification Service starting on port %d...", cfg.Server.Port)

	// Инициализация Kafka publisher
	publisher := infrastructure.NewKafkaPublisher(cfg.Kafka.Brokers)
	defer publisher.Close()

	// Инициализация SendGrid клиента
	sgClient := sendgrid.NewSendClient(cfg.SendGrid.APIKey)
	sendgridNotifier := infrastructure.NewSendgridNotifier(sgClient, cfg.SendGrid.SenderName, cfg.SendGrid.SenderEmail)

	// Инициализация шаблонов (заглушка)
	templateRepo := domain.NewTemplateRepository()

	// Сервис уведомлений (TODO: использовать notificationService в обработчиках событий)
	notificationService := notifier.NewService(sendgridNotifier, templateRepo, publisher)

	// Запуск Kafka consumer
	topics := []string{"weather.updated", "subscription.created", "subscription.confirmed", "subscription.cancelled"}
	consumer := infrastructure.NewKafkaConsumer(cfg.Kafka.Brokers, topics, func(topic string, message []byte) error {
		log.Printf("[HANDLER] Topic: %s, Message: %s", topic, string(message))
		// Business logic for sending email
		switch topic {
		case "weather.updated":
			// Now expecting email in the event
			var event struct {
				City      string                `json:"city"`
				Metrics   domain.WeatherMetrics `json:"metrics"`
				UpdatedAt string                `json:"updated_at"`
				Email     string                `json:"email"`
			}
			if err := json.Unmarshal(message, &event); err != nil {
				return err
			}
			return notificationService.SendWeatherUpdate("email", event.Email, event.Metrics)
		// Add handling for other topics
		default:
			return nil
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumer.Start(ctx)

	// Healthcheck/log
	log.Println("Notification Service is running. Waiting for events...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Notification Service shutting down...")
	cancel()
}
