package events_test

import (
	"Weather-Forecast-API/internal/events"
	"testing"
	"time"
)

func TestWeatherUpdatedEvent(t *testing.T) {
	metrics := map[string]interface{}{
		"temperature": 25.5,
		"humidity":    60.0,
	}

	event := events.WeatherUpdatedEvent{
		City:      "London",
		Metrics:   metrics,
		UpdatedAt: time.Now(),
	}

	if event.City != "London" {
		t.Errorf("Expected city to be 'London', got %s", event.City)
	}

	if event.Metrics["temperature"] != 25.5 {
		t.Errorf("Expected temperature to be 25.5, got %v", event.Metrics["temperature"])
	}
}

func TestSubscriptionCreatedEvent(t *testing.T) {
	event := events.SubscriptionCreatedEvent{
		SubscriptionID:   123,
		ChannelType:      "email",
		ChannelValue:     "test@example.com",
		City:             "Paris",
		FrequencyMinutes: 60,
		Token:            "test-token",
	}

	if event.SubscriptionID != 123 {
		t.Errorf("Expected subscription ID to be 123, got %d", event.SubscriptionID)
	}

	if event.ChannelType != "email" {
		t.Errorf("Expected channel type to be 'email', got %s", event.ChannelType)
	}
}

func TestSubscriptionConfirmedEvent(t *testing.T) {
	event := events.SubscriptionConfirmedEvent{
		SubscriptionID: 123,
		Token:          "test-token",
		ConfirmedAt:    time.Now(),
	}

	if event.SubscriptionID != 123 {
		t.Errorf("Expected subscription ID to be 123, got %d", event.SubscriptionID)
	}

	if event.Token != "test-token" {
		t.Errorf("Expected token to be 'test-token', got %s", event.Token)
	}
}

func TestNotificationSentEvent(t *testing.T) {
	event := events.NotificationSentEvent{
		NotificationID: "notif-123",
		ChannelType:    "email",
		Recipient:      "user@example.com",
		Status:         "sent",
		SentAt:         time.Now(),
	}

	if event.NotificationID != "notif-123" {
		t.Errorf("Expected notification ID to be 'notif-123', got %s", event.NotificationID)
	}

	if event.Status != "sent" {
		t.Errorf("Expected status to be 'sent', got %s", event.Status)
	}
}
