package subscription

import (
	"Weather-Forecast-API/internal/repository"
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) Subscribe(sub *repository.Subscription) error {
	if err := repository.CreateSubscription(sub); err != nil {
		return err
	}
	return nil
}

func (s *Service) Unsubscribe(sub *repository.Subscription) error {
	if err := repository.UnsubscribeByToken(sub.Token); err != nil {
		return err
	}
	return nil
}

func (s *Service) Confirm(sub *repository.Subscription) error {
	if err := repository.ConfirmByToken(sub.Token); err != nil {
		return err
	}
	return nil
}
