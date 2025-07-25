package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type messageWriterManager interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Publisher struct {
	writer messageWriterManager
}

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
