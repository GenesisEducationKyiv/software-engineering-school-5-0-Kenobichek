package handlers

import (
	"context"
	"log"
	"subscription-service/domain"
	"subscription-service/infrastructure"
	"subscription-service/repository"
)

type commandHandler interface {
	Handle(cmd domain.SubscriptionCommand) error
}

type SubscribeHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *SubscribeHandler) Handle(cmd domain.SubscriptionCommand) error {
	sub := &repository.Subscription{
		ChannelType:      cmd.ChannelType,
		ChannelValue:     cmd.ChannelValue,
		City:             cmd.City,
		FrequencyMinutes: cmd.FrequencyMinutes,
		Token:            cmd.Token,
	}
	if err := h.repo.CreateSubscription(sub); err != nil {
		log.Printf("Failed to create subscription: %v", err)
		return err
	}
	log.Printf("Subscription created: %+v", sub)
	event := domain.SubscriptionEvent{
		EventType:        "subscription.created",
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
		Token:            sub.Token,
	}
	if err := h.publisher.Publish(context.Background(), event); err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
	return nil
}

type ConfirmHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *ConfirmHandler) Handle(cmd domain.SubscriptionCommand) error {
	if err := h.repo.ConfirmByToken(cmd.Token); err != nil {
		log.Printf("Failed to confirm subscription: %v", err)
		return err
	}
	log.Printf("Subscription confirmed: %s", cmd.Token)
	event := domain.SubscriptionEvent{
		EventType: "subscription.confirmed",
		Token:     cmd.Token,
	}
	if err := h.publisher.Publish(context.Background(), event); err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
	return nil
}

type UnsubscribeHandler struct {
	repo      *repository.Repository
	publisher *infrastructure.KafkaPublisher
}

func (h *UnsubscribeHandler) Handle(cmd domain.SubscriptionCommand) error {
	if err := h.repo.UnsubscribeByToken(cmd.Token); err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
		return err
	}
	log.Printf("Unsubscribed: %s", cmd.Token)
	event := domain.SubscriptionEvent{
		EventType: "subscription.cancelled",
		Token:     cmd.Token,
	}
	if err := h.publisher.Publish(context.Background(), event); err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
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
