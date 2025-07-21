package handlers_test

import (
	"context"
	"errors"
	"testing"

	"subscription-service/internal/repository"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateSubscriptionFunc     func(ctx context.Context, sub *repository.Subscription) error
	UnsubscribeByTokenFunc     func(ctx context.Context, token string) error
	GetSubscriptionByTokenFunc func(ctx context.Context, token string) (*repository.Subscription, error)
}

func (m *mockRepo) CreateSubscription(ctx context.Context, sub *repository.Subscription) error {
	return m.CreateSubscriptionFunc(ctx, sub)
}
func (m *mockRepo) UnsubscribeByToken(ctx context.Context, token string) error {
	return m.UnsubscribeByTokenFunc(ctx, token)
}
func (m *mockRepo) GetSubscriptionByToken(ctx context.Context, token string) (*repository.Subscription, error) {
	return m.GetSubscriptionByTokenFunc(ctx, token)
}

func TestMockRepo_CreateSubscription(t *testing.T) {
	repo := &mockRepo{
		CreateSubscriptionFunc: func(ctx context.Context, sub *repository.Subscription) error {
			if sub.ChannelValue == "" {
				return errors.New("empty email")
			}
			return nil
		},
	}
	err := repo.CreateSubscription(context.Background(), &repository.Subscription{
		ChannelType:      "email",
		ChannelValue:     "test@example.com",
		City:             "Poltava",
		FrequencyMinutes: 60,
		Token:            "token123",
	})
	assert.NoError(t, err)
}

func TestMockRepo_GetSubscriptionByToken_NotFound(t *testing.T) {
	repo := &mockRepo{
		GetSubscriptionByTokenFunc: func(ctx context.Context, token string) (*repository.Subscription, error) {
			return nil, errors.New("not found")
		},
	}
	_, err := repo.GetSubscriptionByToken(context.Background(), "badtoken")
	assert.EqualError(t, err, "not found")
}
