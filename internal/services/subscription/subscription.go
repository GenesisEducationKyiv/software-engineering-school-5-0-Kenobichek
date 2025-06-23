package subscription

import (
	"Weather-Forecast-API/internal/repository/subscriptions"
	"time"
)

type subscriptionRepositoryManager interface {
	CreateSubscription(subscription *subscriptions.Info) error
	UnsubscribeByToken(token string) error
	ConfirmByToken(token string) error
	GetSubscriptionByToken(token string) (*subscriptions.Info, error)
	GetDueSubscriptions() []subscriptions.Info
	UpdateNextNotification(id int, next time.Time) error
}

type Service struct {
	repo subscriptionRepositoryManager
}

func NewService(repo subscriptionRepositoryManager) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Subscribe(sub *subscriptions.Info) error {
	if err := s.repo.CreateSubscription(sub); err != nil {
		return err
	}
	return nil
}

func (s *Service) Unsubscribe(sub *subscriptions.Info) error {
	if err := s.repo.UnsubscribeByToken(sub.Token); err != nil {
		return err
	}
	return nil
}

func (s *Service) Confirm(sub *subscriptions.Info) error {
	if err := s.repo.ConfirmByToken(sub.Token); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetSubscriptionByToken(token string) (*subscriptions.Info, error) {
	subscription, err := s.repo.GetSubscriptionByToken(token)
	if err != nil {
		return subscription, err
	}
	return subscription, nil
}

func (s *Service) GetDueSubscriptions() []subscriptions.Info {
	return s.repo.GetDueSubscriptions()
}

func (s *Service) UpdateNextNotification(id int, next time.Time) error {
	return s.repo.UpdateNextNotification(id, next)
}