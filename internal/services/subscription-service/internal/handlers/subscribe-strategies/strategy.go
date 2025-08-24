package subscribestrategies

import (
	"context"
	"subscription-service/internal/domain"
	"subscription-service/internal/repository/subscriptions"
)

type loggerManager interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type subscriptionRepositoryManager interface {
	CreateSubscription(ctx context.Context, sub *subscriptions.Subscription) error
	ConfirmByToken(ctx context.Context, token string) error
	UnsubscribeByToken(ctx context.Context, token string) error
	GetSubscriptionByToken(ctx context.Context, token string) (*subscriptions.Subscription, error)
}

type eventPublisherManager interface {
	PublishWithTopic(ctx context.Context, topic string, event interface{}) error
}

type CommandStrategy interface {
	Execute(ctx context.Context, cmd domain.SubscriptionCommand) error
}
