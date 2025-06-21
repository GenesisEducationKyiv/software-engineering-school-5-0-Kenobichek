package subscription

import (
	"Weather-Forecast-API/internal/repository"
	"errors"
	"testing"
)

// Override repository functions for testing
var (
	originalCreateSubscription = repository.CreateSubscription
	originalUnsubscribeByToken = repository.UnsubscribeByToken
	originalConfirmByToken     = repository.ConfirmByToken
)

func TestService_Subscribe(t *testing.T) {
	tests := []struct {
		name          string
		subscription  *repository.Subscription
		repositoryErr error
		expectedErr   string
	}{
		{
			name: "successful subscription",
			subscription: &repository.Subscription{
				ChannelType:      "email",
				ChannelValue:     "test@example.com",
				City:             "London",
				FrequencyMinutes: 1440,
				Token:            "test-token",
			},
			repositoryErr: nil,
			expectedErr:   "",
		},
		{
			name: "repository error",
			subscription: &repository.Subscription{
				ChannelType:      "email",
				ChannelValue:     "test@example.com",
				City:             "London",
				FrequencyMinutes: 1440,
				Token:            "test-token",
			},
			repositoryErr: errors.New("already subscribed"),
			expectedErr:   "already subscribed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock repository function
			repository.CreateSubscription = func(sub *repository.Subscription) error {
				return tt.repositoryErr
			}

			service := NewService()
			err := service.Subscribe(tt.subscription)

			if tt.expectedErr == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.expectedErr != "" && (err == nil || err.Error() != tt.expectedErr) {
				t.Errorf("expected error %s, got %v", tt.expectedErr, err)
			}

			// Restore original function
			repository.CreateSubscription = originalCreateSubscription
		})
	}
}

func TestService_Unsubscribe(t *testing.T) {
	tests := []struct {
		name          string
		subscription  *repository.Subscription
		repositoryErr error
		expectedErr   string
	}{
		{
			name: "successful unsubscribe",
			subscription: &repository.Subscription{
				Token: "valid-token",
			},
			repositoryErr: nil,
			expectedErr:   "",
		},
		{
			name: "repository error",
			subscription: &repository.Subscription{
				Token: "invalid-token",
			},
			repositoryErr: errors.New("not found"),
			expectedErr:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock repository function
			repository.UnsubscribeByToken = func(token string) error {
				return tt.repositoryErr
			}

			service := NewService()
			err := service.Unsubscribe(tt.subscription)

			if tt.expectedErr == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.expectedErr != "" && (err == nil || err.Error() != tt.expectedErr) {
				t.Errorf("expected error %s, got %v", tt.expectedErr, err)
			}

			// Restore original function
			repository.UnsubscribeByToken = originalUnsubscribeByToken
		})
	}
}

func TestService_Confirm(t *testing.T) {
	tests := []struct {
		name          string
		subscription  *repository.Subscription
		repositoryErr error
		expectedErr   string
	}{
		{
			name: "successful confirm",
			subscription: &repository.Subscription{
				Token: "valid-token",
			},
			repositoryErr: nil,
			expectedErr:   "",
		},
		{
			name: "repository error",
			subscription: &repository.Subscription{
				Token: "invalid-token",
			},
			repositoryErr: errors.New("not found"),
			expectedErr:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock repository function
			repository.ConfirmByToken = func(token string) error {
				return tt.repositoryErr
			}

			service := NewService()
			err := service.Confirm(tt.subscription)

			if tt.expectedErr == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.expectedErr != "" && (err == nil || err.Error() != tt.expectedErr) {
				t.Errorf("expected error %s, got %v", tt.expectedErr, err)
			}

			// Restore original function
			repository.ConfirmByToken = originalConfirmByToken
		})
	}
}