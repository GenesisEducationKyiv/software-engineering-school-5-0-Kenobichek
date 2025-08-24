package subscribestrategies

import (
	"context"
	"fmt"
	"subscription-service/internal/domain"
)

type UnsubscribeStrategy struct {
	repo      subscriptionRepositoryManager
	publisher eventPublisherManager
	logger    loggerManager
}

func (u *UnsubscribeStrategy) Execute(ctx context.Context, cmd domain.SubscriptionCommand) error {
	sub, err := u.repo.GetSubscriptionByToken(ctx, cmd.Token)
	if err != nil {
		u.logger.Errorf("Failed to get subscription by token: %v", err)
		return fmt.Errorf("failed to get subscription by token '%s': %w", cmd.Token, err)
	}

	u.logger.Infof("Handling unsubscribe command: %+v", cmd)
	if err := u.repo.UnsubscribeByToken(ctx, cmd.Token); err != nil {
		u.logger.Errorf("Failed to unsubscribe: %v", err)
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	u.logger.Infof("Unsubscribed: %s", cmd.Token)

	event := domain.SubscriptionEvent{
		EventType:        "subscription.cancelled",
		Token:            cmd.Token,
		ChannelType:      sub.ChannelType,
		ChannelValue:     sub.ChannelValue,
		City:             sub.City,
		FrequencyMinutes: sub.FrequencyMinutes,
	}
	u.logger.Infof("Publishing event: %+v", event)
	if err := u.publisher.PublishWithTopic(ctx, "subscription.cancelled", event); err != nil {
		u.logger.Errorf("Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish cancellation event: %w", err)
	}
	u.logger.Infof("Unsubscribe command handled successfully for token=%s", cmd.Token)
	return nil
}
