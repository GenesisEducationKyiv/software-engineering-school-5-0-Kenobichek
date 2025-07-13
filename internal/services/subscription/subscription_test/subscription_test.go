package subscription_test

import (
	"Weather-Forecast-API/internal/events"
	"Weather-Forecast-API/internal/services/subscription"
	"errors"
	"testing"
	"time"

	"Weather-Forecast-API/internal/repository/subscriptions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSubscriptionRepository struct {
	mock.Mock
}

func (m *mockSubscriptionRepository) CreateSubscription(subscription *subscriptions.Info) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *mockSubscriptionRepository) UnsubscribeByToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockSubscriptionRepository) ConfirmByToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockSubscriptionRepository) GetSubscriptionByToken(token string) (*subscriptions.Info, error) {
	args := m.Called(token)
	return args.Get(0).(*subscriptions.Info), args.Error(1)
}

func (m *mockSubscriptionRepository) GetDueSubscriptions() []subscriptions.Info {
	args := m.Called()
	return args.Get(0).([]subscriptions.Info)
}

func (m *mockSubscriptionRepository) UpdateNextNotification(id int, next time.Time) error {
	args := m.Called(id, next)
	return args.Error(0)
}

type mockEventPublisher struct{}

func (m *mockEventPublisher) PublishWeatherUpdated(event events.WeatherUpdatedEvent) error {
	return nil
}
func (m *mockEventPublisher) PublishSubscriptionCreated(event events.SubscriptionCreatedEvent) error {
	return nil
}
func (m *mockEventPublisher) PublishSubscriptionConfirmed(event events.SubscriptionConfirmedEvent) error {
	return nil
}
func (m *mockEventPublisher) PublishSubscriptionCancelled(event events.SubscriptionCancelledEvent) error {
	return nil
}

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name       string
		input      *subscriptions.Info
		mockResult error
		wantError  bool
	}{
		{"valid subscription", &subscriptions.Info{ChannelValue: "test@example.com"}, nil, false},
		{"repository error", &subscriptions.Info{}, errors.New("db error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockSubscriptionRepository)
			mockRepo.On("CreateSubscription", tt.input).Return(tt.mockResult)

			mockPublisher := &mockEventPublisher{}
			service := subscription.NewService(mockRepo, mockPublisher)
			err := service.Subscribe(tt.input)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertCalled(t, "CreateSubscription", tt.input)
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		mockResult error
		wantError  bool
	}{
		{"valid token", "valid-token", nil, false},
		{"repository error", "invalid-token", errors.New("db error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockSubscriptionRepository)
			mockRepo.On("UnsubscribeByToken", tt.input).Return(tt.mockResult)

			mockPublisher := &mockEventPublisher{}
			service := subscription.NewService(mockRepo, mockPublisher)
			err := service.Unsubscribe(&subscriptions.Info{Token: tt.input})

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertCalled(t, "UnsubscribeByToken", tt.input)
		})
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		mockResult error
		wantError  bool
	}{
		{"valid token", "valid-token", nil, false},
		{"repository error", "invalid-token", errors.New("db error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockSubscriptionRepository)
			mockRepo.On("ConfirmByToken", tt.input).Return(tt.mockResult)

			mockPublisher := &mockEventPublisher{}
			service := subscription.NewService(mockRepo, mockPublisher)
			err := service.Confirm(&subscriptions.Info{Token: tt.input})

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertCalled(t, "ConfirmByToken", tt.input)
		})
	}
}

func TestGetSubscriptionByToken(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		mockResult    *subscriptions.Info
		mockError     error
		expected      *subscriptions.Info
		expectedError bool
	}{
		{"valid token", "valid-token", &subscriptions.Info{ChannelValue: "test@example.com"},
			nil, &subscriptions.Info{ChannelValue: "test@example.com"}, false},
		{"repository error", "invalid-token", nil, errors.New("db error"), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockSubscriptionRepository)
			mockRepo.On("GetSubscriptionByToken", tt.input).Return(tt.mockResult, tt.mockError)

			mockPublisher := &mockEventPublisher{}
			service := subscription.NewService(mockRepo, mockPublisher)
			result, err := service.GetSubscriptionByToken(tt.input)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertCalled(t, "GetSubscriptionByToken", tt.input)
		})
	}
}

func TestGetDueSubscriptions(t *testing.T) {
	mockResult := []subscriptions.Info{
		{ChannelValue: "test1@example.com"},
		{ChannelValue: "test2@example.com"},
	}

	mockRepo := new(mockSubscriptionRepository)
	mockRepo.On("GetDueSubscriptions").Return(mockResult)

	mockPublisher := &mockEventPublisher{}
	service := subscription.NewService(mockRepo, mockPublisher)
	result := service.GetDueSubscriptions()

	assert.Equal(t, mockResult, result)

	mockRepo.AssertCalled(t, "GetDueSubscriptions")
}

func TestUpdateNextNotification(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		next       time.Time
		mockResult error
		wantError  bool
	}{
		{"valid update", 1, time.Now().Add(24 * time.Hour), nil, false},
		{"repository error", 2, time.Now().Add(24 * time.Hour), errors.New("db error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockSubscriptionRepository)
			mockRepo.On("UpdateNextNotification", tt.id, tt.next).Return(tt.mockResult)

			mockPublisher := &mockEventPublisher{}
			service := subscription.NewService(mockRepo, mockPublisher)
			err := service.UpdateNextNotification(tt.id, tt.next)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertCalled(t, "UpdateNextNotification", tt.id, tt.next)
		})
	}
}
