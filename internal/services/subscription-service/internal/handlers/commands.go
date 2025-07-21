package handlers

import (
	"context"
	"log"
	"subscription-service/internal/domain"
	"subscription-service/internal/infrastructure"
	"subscription-service/internal/repository"

	"github.com/google/uuid"
)

type SubscribeHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *SubscribeHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	log.Printf("[SubscribeHandler] Handling subscribe command: %+v", cmd)
	sub := &repository.Subscription{
		ChannelType:      cmd.ChannelType,
		ChannelValue:     cmd.ChannelValue,
		City:             cmd.City,
		FrequencyMinutes: cmd.FrequencyMinutes,
		Token:            uuid.NewString(),
	}

	if err := h.repo.CreateSubscription(ctx, sub); err != nil {
		log.Printf("[SubscribeHandler] Failed to create subscription: %v", err)
		return err
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
	}
	log.Printf("[SubscribeHandler] Subscribe command handled successfully for token=%s", sub.Token)
	return nil
}

type ConfirmHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *ConfirmHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {
	log.Printf("[ConfirmHandler] Handling confirm command: %+v", cmd)
	if err := h.repo.ConfirmByToken(ctx, cmd.Token); err != nil {
		log.Printf("[ConfirmHandler] Failed to confirm subscription: %v", err)
		return err
	}
	log.Printf("[ConfirmHandler] Confirm command handled successfully for token=%s", cmd.Token)
	return nil
}

type UnsubscribeHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *UnsubscribeHandler) Handle(ctx context.Context, cmd domain.SubscriptionCommand) error {

	sub, err := h.repo.GetSubscriptionByToken(ctx, cmd.Token)
	if err != nil {
		log.Printf("[UnsubscribeHandler] Failed to get subscription by token: %v", err)
		return err
	}

	log.Printf("[UnsubscribeHandler] Handling unsubscribe command: %+v", cmd)
	if err := h.repo.UnsubscribeByToken(ctx, cmd.Token); err != nil {
		log.Printf("[UnsubscribeHandler] Failed to unsubscribe: %v", err)
		return err
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
	}
	log.Printf("[UnsubscribeHandler] Unsubscribe command handled successfully for token=%s", cmd.Token)
	return nil
}

type commandHandler interface {
	Handle(ctx context.Context, cmd domain.SubscriptionCommand) error
}

type Dispatcher struct {
	handlers map[string]commandHandler
}

func NewDispatcher(repo *repository.Repository, publisher *infrastructure.KafkaPublisher) *Dispatcher {
	return &Dispatcher{
		handlers: map[string]commandHandler{
			"subscribe":   &SubscribeHandler{repo: repo, publisher: publisher},
			"confirm":     &ConfirmHandler{repo: repo, publisher: publisher},
			"unsubscribe": &UnsubscribeHandler{repo: repo, publisher: publisher},
		},
	}
}

func (d *Dispatcher) Handle(cmd domain.SubscriptionCommand) error {
	log.Printf("Received command: %+v", cmd)
	h, ok := d.handlers[cmd.Command]
	if !ok {
		log.Printf("Unknown command: %s", cmd.Command)
		return nil
	}
	return h.Handle(context.Background(), cmd)
}
