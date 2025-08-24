package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"subscription-service/internal/domain"
	subscribestrategies "subscription-service/internal/handlers/subscribe-strategies"

	"github.com/segmentio/kafka-go"
)

const (
	delay                   = 200 * time.Millisecond
	maxRetryDelay           = 30 * time.Second
	maxHandlerRetryAttempts = 5
	maxByteLimit            = 10 * 1024
	minByteLimit            = 10 * 1024
	commitInterval          = 0
)

type loggerManager interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type StrategySelector func(cmd string) (subscribestrategies.CommandStrategy, error)

type KafkaConsumer struct {
	brokers   []string
	topics    []string
	groupID   string
	logger    loggerManager
	selectStrategy StrategySelector
}

func NewKafkaConsumer(
	brokers, topics []string,
	groupID string,
	logger loggerManager,
	selector StrategySelector,
) *KafkaConsumer {
	return &KafkaConsumer{
		brokers:   brokers,
		topics:    topics,
		groupID:   groupID,
		logger:    logger,
		selectStrategy: selector,
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
		Brokers:        c.brokers,
		Topic:          topic,
		GroupID:        c.groupID,
		MinBytes:       minByteLimit,
		MaxBytes:       maxByteLimit,
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

		if err := c.processWithRetry(ctx, topic, m.Value); err != nil {
			c.logger.Infof("handler failed for topic %s: %v", topic, err)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			c.logger.Errorf("failed to commit message for topic %s: %v", topic, err)
			return err
		}
	}
}

func (c *KafkaConsumer) processWithRetry(ctx context.Context, topic string, msg []byte) error {
	retryDelay := delay
	var lastErr error

	for i := 0; i < maxHandlerRetryAttempts; i++ {
		var cmd domain.SubscriptionCommand
		if err := json.Unmarshal(msg, &cmd); err != nil {
			c.logger.Errorf("failed to unmarshal message for topic %s: %v", topic, err)
			lastErr = err
			break
		}

		strategy, err := c.selectStrategy(cmd.Command)
		if err != nil {
			c.logger.Errorf("failed to get strategy for command %s: %v", cmd.Command, err)
			lastErr = err
			break
		}

		err = strategy.Execute(ctx, cmd)
		if err == nil {
			return nil
		}
		lastErr = err
		c.logger.Errorf("strategy execution error (attempt %d/%d) for topic %s: %v", i+1, maxHandlerRetryAttempts, topic, err)

		select {
		case <-time.After(retryDelay):
			retryDelay *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if lastErr == nil {
		lastErr = errors.New("unknown error")
	}
	return errors.New("max handler retry attempts reached: " + lastErr.Error())
}
