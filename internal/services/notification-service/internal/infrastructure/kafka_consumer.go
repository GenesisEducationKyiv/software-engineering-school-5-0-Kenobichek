package infrastructure

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type EventHandler func(topic string, message []byte) error

type KafkaConsumer struct {
	brokers []string
	topics  []string
	groupID string
	handler EventHandler
}

func NewKafkaConsumer(brokers, topics []string, groupID string, handler EventHandler) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topics:  topics,
		groupID: groupID,
		handler: handler,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(c.topics))
	
	go func() {
		wg.Wait()
		close(done)
	}()
	
	for _, topic := range c.topics {
		go func(topic string) {
			defer wg.Done()
			c.consumeTopic(ctx, topic)
		}(topic)
	}
	return done
}
	
func (c *KafkaConsumer) consumeTopic(ctx context.Context, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		Topic:   topic,
		GroupID: c.groupID,
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("failed to close kafka reader: %v", err)
		}
	}()
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
