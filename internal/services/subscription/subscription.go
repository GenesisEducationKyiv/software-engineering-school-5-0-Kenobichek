package subscription

import (
	"Weather-Forecast-API/internal/repository"
)

type SubscriptionService interface {
	Subscribe(sub *repository.Subscription) error
	Unsubscribe(sub *repository.Subscription) error
	Confirm(sub *repository.Subscription) error
}

func NewSubscriptionService() *Subscription {
	return &Subscription{}
}

type Subscription struct{}

func (s *Subscription) Subscribe(sub *repository.Subscription) error {
	if err := repository.CreateSubscription(sub); err != nil {
		return err
	}
	return nil
}

func (s *Subscription) Unsubscribe(sub *repository.Subscription) error {
	if err := repository.UnsubscribeByToken(sub.Token); err != nil {
		return err
	}
	return nil
}

func (s *Subscription) Confirm(sub *repository.Subscription) error {
	if err := repository.ConfirmByToken(sub.Token); err != nil {
		return err
	}
	return nil
}
