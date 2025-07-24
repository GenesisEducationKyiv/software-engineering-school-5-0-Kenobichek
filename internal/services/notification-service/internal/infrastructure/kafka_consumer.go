package infrastructure

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventHandler func(ctx context.Context, topic string, message []byte) error

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

	for _, topic := range c.topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			c.consumeTopicWithRetries(ctx, topic)
		}(topic)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

func (c *KafkaConsumer) consumeTopicWithRetries(ctx context.Context, topic string) {
	const maxRetryDelay = 30 * time.Second
	retryDelay := 500 * time.Millisecond

	for {
		err := c.consumeTopic(ctx, topic)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			log.Printf("[WARN] consumer for topic %s failed: %v, retrying in %v", topic, err, retryDelay)
			select {
			case <-time.After(retryDelay):
				retryDelay *= 2
				if retryDelay > maxRetryDelay {
					retryDelay = maxRetryDelay
				}
			case <-ctx.Done():
				log.Printf("[INFO] context cancelled during retry wait for topic %s", topic)
				return
			}
			continue
		}
		return
	}
}

func (c *KafkaConsumer) consumeTopic(ctx context.Context, topic string) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		Topic:    topic,
		GroupID:  c.groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("[ERROR] failed to close kafka reader for topic %s: %v", topic, err)
		}
	}()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Printf("[INFO] consumer for topic %s stopped due to context: %v", topic, err)
				return err
			}

			log.Printf("[ERROR] kafka fetch error for topic %s: %v", topic, err)
			return err
		}

		log.Printf("[KAFKA CONSUMER] received event from topic %s: %s", topic, string(m.Value))

		if c.handler != nil {
			if err := c.processWithRetry(ctx, topic, m.Value); err != nil {
				log.Printf("[ERROR] handler failed for topic %s: %v", topic, err)
				continue
			}
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Printf("[ERROR] failed to commit message for topic %s: %v", topic, err)
			return err
		}
	}
}

func (c *KafkaConsumer) processWithRetry(ctx context.Context, topic string, msg []byte) error {
	const maxAttempts = 5
	delay := 200 * time.Millisecond
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		err := c.handler(ctx, topic, msg)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("[WARN] handler error (attempt %d/%d) for topic %s: %v", i+1, maxAttempts, topic, err)

		select {
		case <-time.After(delay):
			delay *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.New("max handler retry attempts reached: " + lastErr.Error())
}
