package handlers

import (
	"context"
	"fmt"
	"log"

	"subscription-service/internal/domain"
	"subscription-service/internal/repository/subscriptions"

	"github.com/google/uuid"
)

const (
	subscribeCommand = "subscribe"
	confirmCommand   = "confirm"
	unsubscribeCommand = "unsubscribe"
)

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
}

func (h *SubscribeHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	log.Printf("[SubscribeHandler] Handling subscribe command: %+v", cmd)
	sub := &subscriptions.Subscription{
		ChannelType:      cmd.ChannelType,
		ChannelValue:     cmd.ChannelValue,
		City:             cmd.City,
		FrequencyMinutes: cmd.FrequencyMinutes,
		Token:            uuid.NewString(),
	}

	if err := h.repo.CreateSubscription(ctx, sub); err != nil {
		log.Printf("[SubscribeHandler] Failed to create subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	log.Printf("[SubscribeHandler] Subscription created: %+v", sub)
	event := domain.SubscriptionEvent{
		EventType:        "subscription.confirmed",
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
		Token:            sub.Token,
	}
	log.Printf("[SubscribeHandler] Publishing event: %+v", event)
	if err := h.publisher.PublishWithTopic(ctx, "subscription.confirmed", event); err != nil {
		log.Printf("[SubscribeHandler] Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish confirmation event: %w", err)
	}
	log.Printf("[SubscribeHandler] Subscribe command handled successfully for token=%s", sub.Token)
	return nil
}

type ConfirmHandler struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
}

func (h *ConfirmHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	log.Printf("[ConfirmHandler] Handling confirm command: %+v", cmd)
	if err := h.repo.ConfirmByToken(ctx, cmd.Token); err != nil {
		log.Printf("[ConfirmHandler] Failed to confirm subscription: %v", err)
		return fmt.Errorf("failed to confirm subscription: %w", err)
	}
	log.Printf("[ConfirmHandler] Confirm command handled successfully for token=%s", cmd.Token)
	return nil
}

type UnsubscribeHandler struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
}

func (h *UnsubscribeHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {

	sub, err := h.repo.GetSubscriptionByToken(ctx, cmd.Token)
	if err != nil {
		log.Printf("[UnsubscribeHandler] Failed to get subscription by token: %v", err)
		return fmt.Errorf("failed to get subscription by token '%s': %w", cmd.Token, err)
	}

	log.Printf("[UnsubscribeHandler] Handling unsubscribe command: %+v", cmd)
	if err := h.repo.UnsubscribeByToken(ctx, cmd.Token); err != nil {
		log.Printf("[UnsubscribeHandler] Failed to unsubscribe: %v", err)
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	log.Printf("[UnsubscribeHandler] Unsubscribed: %s", cmd.Token)
	event := domain.SubscriptionEvent{
		EventType:        "subscription.cancelled",
		Token:            cmd.Token,
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
	}
	log.Printf("[UnsubscribeHandler] Publishing event: %+v", event)
	if err := h.publisher.PublishWithTopic(ctx, "subscription.cancelled", event); err != nil {
		log.Printf("[UnsubscribeHandler] Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish cancellation event: %w", err)
	}
	log.Printf("[UnsubscribeHandler] Unsubscribe command handled successfully for token=%s", cmd.Token)
	return nil
}

type commandHandler interface {
	Handle(ctx context.Context, cmd domain.SubscriptionCommand) error
}

type dispatcher struct {
	handlers map[string]commandHandler
}

func NewDispatcher(repo subscriptionRepositoryManager, publisher eventPublisherManager) *dispatcher {
	return &dispatcher{
		handlers: map[string]commandHandler{
			subscribeCommand:   &SubscribeHandler{repo: repo, publisher: publisher},
			confirmCommand:     &ConfirmHandler{repo: repo, publisher: publisher},
			unsubscribeCommand: &UnsubscribeHandler{repo: repo, publisher: publisher},
		},
	}
}

func (d *dispatcher) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	log.Printf("Received command: %+v", cmd)
	h, ok := d.handlers[cmd.Command]
	if !ok {
		log.Printf("Unknown command: %s", cmd.Command)
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
	return h.Handle(ctx, cmd)
}
