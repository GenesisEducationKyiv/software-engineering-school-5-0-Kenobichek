package subscribestrategies

import "fmt"

const (
	subscribeCommand = "subscribe"
	confirmCommand   = "confirm"
	unsubscribeCommand = "unsubscribe"
)

func StrategyFactory(
	cmd string,
	repo subscriptionRepositoryManager,
	publisher eventPublisherManager,
	logger loggerManager,
) (CommandStrategy, error) {
	switch cmd {
	case subscribeCommand:
		return &SubscribeStrategy{
			repo:      repo,
			publisher: publisher,
			logger:    logger,
		}, nil
	case confirmCommand:
		return &ConfirmStrategy{
			repo:   repo,
			logger: logger,
		}, nil
	case unsubscribeCommand:
		return &UnsubscribeStrategy{
			repo:      repo,
			publisher: publisher,
			logger:    logger,
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}
