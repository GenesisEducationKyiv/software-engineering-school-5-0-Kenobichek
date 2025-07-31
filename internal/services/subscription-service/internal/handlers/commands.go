package handlers

import (
	"context"
	"fmt"

	"subscription-service/internal/domain"
	"subscription-service/internal/repository/subscriptions"

	"github.com/google/uuid"
)

const (
	subscribeCommand = "subscribe"
	confirmCommand   = "confirm"
	unsubscribeCommand = "unsubscribe"
)

type loggerManager interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Sync() error
}

type subscriptionRepositoryManager interface {
    CreateSubscription(ctx context.Context, sub *subscriptions.Subscription) error
    ConfirmByToken(ctx context.Context, token string) error
    UnsubscribeByToken(ctx context.Context, token string) error
    GetSubscriptionByToken(ctx context.Context, token string) (*subscriptions.Subscription, error)
}

type eventPublisherManager interface {
    PublishWithTopic(ctx context.Context, topic string, event interface{}) error
}

type SubscribeHandler struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
	logger loggerManager
}

func (h *SubscribeHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	h.logger.Info("Handling subscribe command: %+v", cmd)
	sub := &subscriptions.Subscription{
		ChannelType:      cmd.ChannelType,
		ChannelValue:     cmd.ChannelValue,
		City:             cmd.City,
		FrequencyMinutes: cmd.FrequencyMinutes,
		Token:            uuid.NewString(),
	}

	if err := h.repo.CreateSubscription(ctx, sub); err != nil {
		h.logger.Error("Failed to create subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	h.logger.Info("Subscription created: %+v", sub)
	event := domain.SubscriptionEvent{
		EventType:        "subscription.confirmed",
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
		Token:            sub.Token,
	}
	h.logger.Info("Publishing event: %+v", event)
	if err := h.publisher.PublishWithTopic(ctx, "subscription.confirmed", event); err != nil {
		h.logger.Error("Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish confirmation event: %w", err)
	}
	h.logger.Info("Subscribe command handled successfully for token=%s", sub.Token)
	return nil
}

type ConfirmHandler struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
	logger loggerManager
}

func (h *ConfirmHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	h.logger.Info("Handling confirm command: %+v", cmd)
	if err := h.repo.ConfirmByToken(ctx, cmd.Token); err != nil {
		h.logger.Error("Failed to confirm subscription: %v", err)
		return fmt.Errorf("failed to confirm subscription: %w", err)
	}
	h.logger.Info("Confirm command handled successfully for token=%s", cmd.Token)
	return nil
}

type UnsubscribeHandler struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
	logger loggerManager
}

func (h *UnsubscribeHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {

	sub, err := h.repo.GetSubscriptionByToken(ctx, cmd.Token)
	if err != nil {
		h.logger.Error("Failed to get subscription by token: %v", err)
		return fmt.Errorf("failed to get subscription by token '%s': %w", cmd.Token, err)
	}

	h.logger.Info("Handling unsubscribe command: %+v", cmd)
	if err := h.repo.UnsubscribeByToken(ctx, cmd.Token); err != nil {
		h.logger.Error("Failed to unsubscribe: %v", err)
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	h.logger.Info("Unsubscribed: %s", cmd.Token)
	event := domain.SubscriptionEvent{
		EventType:        "subscription.cancelled",
		Token:            cmd.Token,
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
	}
	h.logger.Info("Publishing event: %+v", event)
	if err := h.publisher.PublishWithTopic(ctx, "subscription.cancelled", event); err != nil {
		h.logger.Error("Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish cancellation event: %w", err)
	}
	h.logger.Info("Unsubscribe command handled successfully for token=%s", cmd.Token)
	return nil
}

type commandHandler interface {
	Handle(ctx context.Context, cmd domain.SubscriptionCommand) error
}

type dispatcher struct {
	handlers map[string]commandHandler
	logger loggerManager
}

func NewDispatcher(repo subscriptionRepositoryManager, publisher eventPublisherManager, logger loggerManager) *dispatcher {
	return &dispatcher{
		handlers: map[string]commandHandler{
			subscribeCommand:   &SubscribeHandler{repo: repo, publisher: publisher, logger: logger},
			confirmCommand:     &ConfirmHandler{repo: repo, publisher: publisher, logger: logger},
			unsubscribeCommand: &UnsubscribeHandler{repo: repo, publisher: publisher, logger: logger},
		},
		logger: logger,
	}
}

func (d *dispatcher) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	d.logger.Info("Received command: %+v", cmd)
	h, ok := d.handlers[cmd.Command]
	if !ok {
		d.logger.Error("Unknown command: %s", cmd.Command)
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
	return h.Handle(ctx, cmd)
}
