package handlers

import (
	"context"
	"log"
	"subscription-service/domain"
	"subscription-service/infrastructure"
	"subscription-service/repository"

	"github.com/google/uuid"
)

type commandHandler interface {
	Handle(cmd domain.SubscriptionCommand) error
}

type SubscribeHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *SubscribeHandler) Handle(cmd domain.SubscriptionCommand) error {
	log.Printf("[SubscribeHandler] Handling subscribe command: %+v", cmd)
	sub := &repository.Subscription{
		ChannelType:      cmd.ChannelType,
		ChannelValue:     cmd.ChannelValue,
		City:             cmd.City,
		FrequencyMinutes: cmd.FrequencyMinutes,
		Token:            uuid.NewString(),
	}

	if err := h.repo.CreateSubscription(sub); err != nil {
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
	if err := h.publisher.PublishWithTopic(context.Background(), "subscription.confirmed", event); err != nil {
		log.Printf("[SubscribeHandler] Failed to publish event: %v", err)
	}
	log.Printf("[SubscribeHandler] Subscribe command handled successfully for token=%s", sub.Token)
	return nil
}

type ConfirmHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *ConfirmHandler) Handle(cmd domain.SubscriptionCommand) error {
	log.Printf("[ConfirmHandler] Handling confirm command: %+v", cmd)
	if err := h.repo.ConfirmByToken(cmd.Token); err != nil {
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

func (h *UnsubscribeHandler) Handle(cmd domain.SubscriptionCommand) error {
	log.Printf("[UnsubscribeHandler] Handling unsubscribe command: %+v", cmd)
	if err := h.repo.UnsubscribeByToken(cmd.Token); err != nil {
		log.Printf("[UnsubscribeHandler] Failed to unsubscribe: %v", err)
		return err
	}
	log.Printf("[UnsubscribeHandler] Unsubscribed: %s", cmd.Token)
	event := domain.SubscriptionEvent{
		EventType: "subscription.cancelled",
		Token:     cmd.Token,
	}
	log.Printf("[UnsubscribeHandler] Publishing event: %+v", event)
	if err := h.publisher.PublishWithTopic(context.Background(), "subscription.cancelled", event); err != nil {
		log.Printf("[UnsubscribeHandler] Failed to publish event: %v", err)
	}
	log.Printf("[UnsubscribeHandler] Unsubscribe command handled successfully for token=%s", cmd.Token)
	return nil
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
	return h.Handle(cmd)
}
