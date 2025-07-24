package handlers

import (
	"encoding/json"
	"fmt"
	"notification-service/internal/domain"
	"notification-service/internal/notifier"
)

func parseEvent[T any](message []byte) (T, error) {
	var event T
	if err := json.Unmarshal(message, &event); err != nil {
		return event, fmt.Errorf("failed to parse event: %w", err)
	}
	return event, nil
}

type WeatherUpdateEvent struct {
	Metrics   domain.WeatherMetrics	`json:"metrics"`
	UpdatedAt int64					`json:"updated_at"`
	Email     string				`json:"channel_value"`
}

type SubscriptionConfirmedEvent struct {
	Email string `json:"channel_value"`
	Token string `json:"token"`
}

type SubscriptionCancelledEvent struct {
	Email string `json:"channel_value"`
	City  string `json:"city"`
}

type WeatherUpdateHandler struct {
	notificationService *notifier.Service
}

func NewWeatherUpdateHandler(service *notifier.Service) *WeatherUpdateHandler {
	return &WeatherUpdateHandler{
		notificationService: service,
	}
}

func (h *WeatherUpdateHandler) Handle(message []byte) error {
	event, err := parseEvent[WeatherUpdateEvent](message)
	if err != nil {
		return err
	}
	return h.notificationService.SendWeatherUpdate("email", event.Email, event.Metrics)
}

type SubscriptionConfirmedHandler struct {
	notificationService *notifier.Service
}

func NewSubscriptionConfirmedHandler(service *notifier.Service) *SubscriptionConfirmedHandler {
	return &SubscriptionConfirmedHandler{
		notificationService: service,
	}
}

func (h *SubscriptionConfirmedHandler) Handle(message []byte) error {
	event, err := parseEvent[SubscriptionConfirmedEvent](message)
	if err != nil {
		return err
	}
	return h.notificationService.SendConfirmation("email", event.Email, event.Token)
}

type SubscriptionCancelledHandler struct {
	notificationService *notifier.Service
}

func NewSubscriptionCancelledHandler(service *notifier.Service) *SubscriptionCancelledHandler {
	return &SubscriptionCancelledHandler{
		notificationService: service,
	}
}

func (h *SubscriptionCancelledHandler) Handle(message []byte) error {
	event, err := parseEvent[SubscriptionCancelledEvent](message)
	if err != nil {
		return err
	}
	return h.notificationService.SendUnsubscribe("email", event.Email, event.City)
}
