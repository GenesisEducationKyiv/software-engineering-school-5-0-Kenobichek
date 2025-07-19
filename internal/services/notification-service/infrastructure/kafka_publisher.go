package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"notification-service/domain"

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
	return &KafkaPublisher{writer: writer}
}

func (k *KafkaPublisher) PublishNotificationSent(event domain.NotificationSentEvent) error {
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
