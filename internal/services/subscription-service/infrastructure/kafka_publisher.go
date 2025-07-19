package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

// KafkaPublisher публикует события SubscriptionEvent в Kafka
// Используйте один экземпляр на сервис

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event interface{}) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{Value: msg})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
