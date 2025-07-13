package events

import (
	"time"
)

type WeatherUpdatedEvent struct {
	City      string                 `json:"city"`
	Metrics   map[string]interface{} `json:"metrics"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type SubscriptionCreatedEvent struct {
	SubscriptionID   int    `json:"subscription_id"`
	ChannelType      string `json:"channel_type"`
	ChannelValue     string `json:"channel_value"`
	City             string `json:"city"`
	FrequencyMinutes int    `json:"frequency_minutes"`
	Token            string `json:"token"`
}

type SubscriptionConfirmedEvent struct {
	SubscriptionID int       `json:"subscription_id"`
	Token          string    `json:"token"`
	ConfirmedAt    time.Time `json:"confirmed_at"`
}

type SubscriptionCancelledEvent struct {
	SubscriptionID int       `json:"subscription_id"`
	Token          string    `json:"token"`
	CancelledAt    time.Time `json:"cancelled_at"`
}

type NotificationSentEvent struct {
	NotificationID string    `json:"notification_id"`
	ChannelType    string    `json:"channel_type"`
	Recipient      string    `json:"recipient"`
	Status         string    `json:"status"`
	SentAt         time.Time `json:"sent_at"`
}

// EventPublisher defines the interface for publishing events
// type EventPublisher interface {
// 	PublishWeatherUpdated(event WeatherUpdatedEvent) error
// 	PublishSubscriptionCreated(event SubscriptionCreatedEvent) error
// 	PublishSubscriptionConfirmed(event SubscriptionConfirmedEvent) error
// 	PublishSubscriptionCancelled(event SubscriptionCancelledEvent) error
// 	PublishNotificationSent(event NotificationSentEvent) error
// }

// // EventConsumer defines the interface for consuming events
// type EventConsumer interface {
// 	Subscribe(topic string, handler interface{}) error
// 	Start() error
// 	Stop() error
// }
