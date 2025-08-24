package subscribestrategies

import (
	"context"
	"fmt"
	"subscription-service/internal/domain"
	"subscription-service/internal/repository/subscriptions"

	"github.com/google/uuid"
)

type SubscribeStrategy struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
	logger    loggerManager
}

func (s *SubscribeStrategy) Execute(ctx context.Context, cmd domain.SubscriptionCommand) error {
	s.logger.Infof("Handling subscribe command: %+v", cmd)
	sub := &subscriptions.Subscription{
		ChannelType:      cmd.ChannelType,
		ChannelValue:     cmd.ChannelValue,
		City:             cmd.City,
		FrequencyMinutes: cmd.FrequencyMinutes,
		Token:            uuid.NewString(),
	}

	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		s.logger.Errorf("Failed to create subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	s.logger.Infof("Subscription created: %+v", sub)

	event := domain.SubscriptionEvent{
		EventType:        "subscription.confirmed",
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
		Token:            sub.Token,
	}
	s.logger.Infof("Publishing event: %+v", event)
	if err := s.publisher.PublishWithTopic(ctx, "subscription.confirmed", event); err != nil {
		s.logger.Errorf("Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish confirmation event: %w", err)
	}
	s.logger.Infof("Subscribe command handled successfully for token=%s", sub.Token)
	return nil
}
