package subscription

import (
	"Weather-Forecast-API/internal/repository"
)

type Subscription interface {
	Subscribe(sub *repository.Subscription) error
	Unsubscribe(sub *repository.Subscription) error
	Confirm(sub *repository.Subscription) error
}

func NewSubscriptionService() SubscriptionService {
	return SubscriptionService{}
}

type SubscriptionService struct{}

func (s *SubscriptionService) Subscribe(sub *repository.Subscription) error {
	if err := repository.CreateSubscription(sub); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) Unsubscribe(sub *repository.Subscription) error {
	if err := repository.UnsubscribeByToken(sub.Token); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) Confirm(sub *repository.Subscription) error {
	if err := repository.ConfirmByToken(sub.Token); err != nil {
		return err
	}
	return nil
}
