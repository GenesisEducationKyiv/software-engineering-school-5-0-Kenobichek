package service

import (
	"Weather-Forecast-API/internal/models"
	"Weather-Forecast-API/internal/repository"
)

type SubscriptionService interface {
	Subscribe(sub *models.Subscription) error
	Unsubscribe(sub *models.Subscription) error
	Confirm(sub *models.Subscription) error
}

func NewSubscriptionService() SubscriptionService {
	return &subscriptionService{}
}
type subscriptionService struct {
}

func (s *subscriptionService) Subscribe(sub *models.Subscription) error {
	if err := repository.CreateSubscription(sub); err != nil {
		return err
	}
	return nil
}

func (s *subscriptionService) Unsubscribe(sub *models.Subscription) error {
	if err := repository.UnsubscribeByToken(sub.Token); err != nil {
		return err
	}
	return nil
}

func (s *subscriptionService) Confirm(sub *models.Subscription) error {
	if err := repository.ConfirmByToken(sub.Token); err != nil {
		return err
	}
	return nil
}
