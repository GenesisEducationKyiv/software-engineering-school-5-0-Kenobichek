package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	"subscription-service/internal/domain"

	"github.com/segmentio/kafka-go"
)

type CommandHandler func(ctx context.Context, cmd domain.SubscriptionCommand) error

type KafkaConsumer struct {
	brokers []string
	topic   string
	handler CommandHandler
}

func NewKafkaConsumer(brokers []string, topic string, handler CommandHandler) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topic:   topic,
		handler: handler,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		Topic:   c.topic,
		GroupID: "subscription-service",
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Kafka reader close error: %v", err)
		}
	}()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Kafka consumer stopped")
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}
		var cmd domain.SubscriptionCommand
		if err := json.Unmarshal(m.Value, &cmd); err != nil {
			log.Printf("Failed to unmarshal command: %v", err)
			continue
		}
		if err := c.handler(ctx, cmd); err != nil {
			log.Printf("Command handler error: %v", err)
		}
	}
}
