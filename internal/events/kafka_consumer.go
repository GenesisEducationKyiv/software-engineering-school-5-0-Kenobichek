package events

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type EventHandler func([]byte) error

type KafkaConsumer struct {
	reader   *kafka.Reader
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	handlers map[string]EventHandler
}

func NewKafkaConsumer(brokers []string, topic, groupID string, handlers map[string]EventHandler) *KafkaConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &KafkaConsumer{
		reader:   reader,
		ctx:      ctx,
		cancel:   cancel,
		handlers: handlers,
	}
}

func (k *KafkaConsumer) Start() error {
	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for {
			select {
			case <-k.ctx.Done():
				return
			default:
				m, err := k.reader.ReadMessage(k.ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}
				topic := k.reader.Config().Topic
				handler, ok := k.handlers[topic]
				if !ok {
					log.Printf("No handler registered for topic: %s", topic)
					continue
				}
				if err := handler(m.Value); err != nil {
					log.Printf("Error handling event for topic %s: %v", topic, err)
				}
			}
		}
	}()
	log.Printf("Started Kafka consumer for topic: %s, group: %s", k.reader.Config().Topic, k.reader.Config().GroupID)
	return nil
}

func (k *KafkaConsumer) Stop() error {
	k.cancel()
	k.wg.Wait()
	return k.reader.Close()
}

func NewWeatherEventLogger() EventHandler {
	return func(data []byte) error {
		var event WeatherUpdatedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		log.Printf("Received weather update for city: %s", event.City)
		return nil
	}
}

func NewSubscriptionEventLogger() EventHandler {
	return func(data []byte) error {
		var event SubscriptionCreatedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		log.Printf("Received subscription created for city: %s", event.City)
		return nil
	}
}
