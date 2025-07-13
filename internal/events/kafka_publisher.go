package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) *KafkaPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaPublisher{
		writer: writer,
	}
}

func (k *KafkaPublisher) PublishWeatherUpdated(event WeatherUpdatedEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal weather event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "weather.updated",
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish weather event: %w", err)
	}

	log.Printf("Published weather updated event for city: %s", event.City)
	return nil
}

func (k *KafkaPublisher) PublishSubscriptionCreated(event SubscriptionCreatedEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "subscription.created",
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish subscription event: %w", err)
	}

	log.Printf("Published subscription created event for city: %s", event.City)
	return nil
}

func (k *KafkaPublisher) PublishSubscriptionConfirmed(event SubscriptionConfirmedEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "subscription.confirmed",
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish subscription event: %w", err)
	}

	log.Printf("Published subscription confirmed event for token: %s", event.Token)
	return nil
}

func (k *KafkaPublisher) PublishSubscriptionCancelled(event SubscriptionCancelledEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "subscription.cancelled",
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish subscription event: %w", err)
	}

	log.Printf("Published subscription cancelled event for token: %s", event.Token)
	return nil
}

func (k *KafkaPublisher) PublishNotificationSent(event NotificationSentEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal notification event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "notification.sent",
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish notification event: %w", err)
	}

	log.Printf("Published notification sent event for recipient: %s", event.Recipient)
	return nil
}

func (k *KafkaPublisher) Close() error {
	return k.writer.Close()
}
