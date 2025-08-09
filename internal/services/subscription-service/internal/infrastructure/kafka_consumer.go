package infrastructure

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	delay = 200 * time.Millisecond
	maxRetryDelay = 30 * time.Second
	maxHandlerRetryAttempts = 5
	maxByteLimit = 10 * 1024
	minByteLimit = 10 * 1024
	commitInterval = 0
)

type loggerManager interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type EventHandler func(ctx context.Context, topic string, message []byte) error

type KafkaConsumer struct {
	brokers []string
	topics  []string
	groupID string
	handler EventHandler
	logger loggerManager
}

func NewKafkaConsumer(
	brokers, topics []string,
	groupID string,
	handler EventHandler,
	logger loggerManager,
) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topics:  topics,
		groupID: groupID,
		handler: handler,
		logger: logger,
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
	retryDelay := delay

	for {
		err := c.consumeTopic(ctx, topic)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			c.logger.Infof("consumer for topic %s failed: %v, retrying in %v", topic, err, retryDelay)
			select {
			case <-time.After(retryDelay):
				retryDelay *= 2
				if retryDelay > maxRetryDelay {
					retryDelay = maxRetryDelay
				}
			case <-ctx.Done():
				c.logger.Infof("[INFO] context cancelled during retry wait for topic %s", topic)
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
		MinBytes: minByteLimit,
		MaxBytes: maxByteLimit,
		CommitInterval: commitInterval,
	})
	defer func() {
		if err := r.Close(); err != nil {
			c.logger.Errorf("failed to close kafka reader for topic %s: %v", topic, err)
		}
	}()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				c.logger.Infof("consumer for topic %s stopped due to context: %v", topic, err)
				return err
			}

			c.logger.Infof("kafka fetch error for topic %s: %v", topic, err)
			return err
		}

		c.logger.Infof("received event from topic %s, partition %d, offset %d", topic, m.Partition, m.Offset)

		if c.handler != nil {
			if err := c.processWithRetry(ctx, topic, m.Value); err != nil {
				c.logger.Infof("handler failed for topic %s: %v", topic, err)
				continue
			}
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			c.logger.Errorf("failed to commit message for topic %s: %v", topic, err)
			return err
		}
	}
}

func (c *KafkaConsumer) processWithRetry(ctx context.Context, topic string, msg []byte) error {
	delay := delay
	var lastErr error

	for i := 0; i < maxHandlerRetryAttempts; i++ {
		err := c.handler(ctx, topic, msg)
		if err == nil {
			return nil
		}
		lastErr = err
		c.logger.Errorf("handler error (attempt %d/%d) for topic %s: %v", i+1, maxHandlerRetryAttempts, topic, err)

		select {
		case <-time.After(delay):
			delay *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.New("max handler retry attempts reached: " + lastErr.Error())
}
