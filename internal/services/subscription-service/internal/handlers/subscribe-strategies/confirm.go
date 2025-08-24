package subscribestrategies

import (
	"context"
	"fmt"
	"subscription-service/internal/domain"
)

type ConfirmStrategy struct {
	repo   subscriptionRepositoryManager
	logger loggerManager
}

func (c *ConfirmStrategy) Execute(ctx context.Context, cmd domain.SubscriptionCommand) error {
	c.logger.Infof("Handling confirm command: %+v", cmd)
	if err := c.repo.ConfirmByToken(ctx, cmd.Token); err != nil {
		c.logger.Errorf("Failed to confirm subscription: %v", err)
		return fmt.Errorf("failed to confirm subscription: %w", err)
	}
	c.logger.Infof("Confirm command handled successfully for token=%s", cmd.Token)
	return nil
}
