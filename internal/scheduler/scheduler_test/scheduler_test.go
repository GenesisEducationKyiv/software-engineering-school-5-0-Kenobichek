package scheduler_test

import (
	"context"
	"testing"
	"time"

	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/scheduler"
)

// Mock implementations for testing
type mockNotificationManager struct{}

func (m *mockNotificationManager) SendWeatherUpdate(channel string, recipient string, metrics weather.Metrics) error {
	return nil
}

type mockSubscriptionManager struct{}

func (m *mockSubscriptionManager) GetSubscriptionByToken(token string) (*subscriptions.Info, error) {
	return nil, nil
}

func (m *mockSubscriptionManager) GetDueSubscriptions() []subscriptions.Info {
	return []subscriptions.Info{}
}

func (m *mockSubscriptionManager) UpdateNextNotification(id int, next time.Time) error {
	return nil
}

type mockWeatherChainHandler struct{}

func (m *mockWeatherChainHandler) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	return weather.Metrics{}, nil
}

func TestSchedulerStop(t *testing.T) {
	// Create a scheduler with mock dependencies
	scheduler := scheduler.NewScheduler(
		&mockNotificationManager{},
		&mockSubscriptionManager{},
		&mockWeatherChainHandler{},
		5*time.Second,
	)

	// Start the scheduler
	_, err := scheduler.Start()
	if err != nil {
		t.Fatalf("Failed to start scheduler: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop the scheduler
	err = scheduler.Stop()
	if err != nil {
		t.Fatalf("Failed to stop scheduler: %v", err)
	}

	// Verify that the scheduler is stopped by trying to stop it again
	// This should not cause any issues
	err = scheduler.Stop()
	if err != nil {
		t.Fatalf("Failed to stop scheduler again: %v", err)
	}
}

func TestSchedulerStopWithoutStart(t *testing.T) {
	// Create a scheduler with mock dependencies
	scheduler := scheduler.NewScheduler(
		&mockNotificationManager{},
		&mockSubscriptionManager{},
		&mockWeatherChainHandler{},
		5*time.Second,
	)

	// Try to stop without starting - should not cause any issues
	err := scheduler.Stop()
	if err != nil {
		t.Fatalf("Failed to stop scheduler without starting: %v", err)
	}
}
