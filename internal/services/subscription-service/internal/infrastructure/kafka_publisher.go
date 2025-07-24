package infrastructure

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	brokers []string
	writers map[string]*kafka.Writer
	mu      sync.Mutex
}

func NewKafkaPublisher(brokers []string, _ string) *KafkaPublisher {
	return &KafkaPublisher{
		brokers: brokers,
		writers: make(map[string]*kafka.Writer),
	}
}

func (p *KafkaPublisher) getWriter(topic string) *kafka.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()

	if w, ok := p.writers[topic]; ok {
		return w
	}
	w := &kafka.Writer{
		Addr:     kafka.TCP(p.brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	p.writers[topic] = w
	p.mu.Unlock()
	return w
}

func (p *KafkaPublisher) PublishWithTopic(ctx context.Context, topic string, event interface{}) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	writer := p.getWriter(topic)
	return writer.WriteMessages(ctx, kafka.Message{Value: msg})
}

func (p *KafkaPublisher) Publish(ctx context.Context, event interface{}) error {
	return p.PublishWithTopic(ctx, "default", event)
}

func (p *KafkaPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	var firstErr error
	for _, w := range p.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	p.writers = make(map[string]*kafka.Writer)
	return firstErr
}
