package infrastructure

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventHandler func(topic string, message []byte) error

type KafkaConsumer struct {
	brokers []string
	topics  []string
	handler EventHandler
}

func NewKafkaConsumer(brokers, topics []string, handler EventHandler) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topics:  topics,
		handler: handler,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	for _, topic := range c.topics {
		go c.consumeTopic(ctx, topic)
	}
}

func (c *KafkaConsumer) consumeTopic(ctx context.Context, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		Topic:   topic,
		GroupID: "notification-service",
	})
	defer r.Close()
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("Kafka consumer for topic %s stopped", topic)
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}
		log.Printf("Received event from topic %s: %s", topic, string(m.Value))
		if c.handler != nil {
			if err := c.handler(topic, m.Value); err != nil {
				log.Printf("Event handler error: %v", err)
			}
		}
	}
}
