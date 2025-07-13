package subscription

import (
	"Weather-Forecast-API/internal/events"
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

type eventPublisherManagerher interface {
	PublishWeatherUpdated(event events.WeatherUpdatedEvent) error
	PublishSubscriptionCreated(event events.SubscriptionCreatedEvent) error
	PublishSubscriptionConfirmed(event events.SubscriptionConfirmedEvent) error
	PublishSubscriptionCancelled(event events.SubscriptionCancelledEvent) error
}


type Service struct {
	repo           subscriptionRepositoryManager
	eventPublisher eventPublisherManagerher
}

func NewService(repo subscriptionRepositoryManager, eventPublisher eventPublisherManagerher) *Service {
	return &Service{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (s *Service) Subscribe(sub *subscriptions.Info) error {
	if err := s.repo.CreateSubscription(sub); err != nil {
		return err
	}
	if s.eventPublisher != nil {
		event := events.SubscriptionCreatedEvent{
			SubscriptionID:   sub.ID,
			ChannelType:      sub.ChannelType,
			ChannelValue:     sub.ChannelValue,
			City:             sub.City,
			FrequencyMinutes: sub.FrequencyMinutes,
			Token:            sub.Token,
		}
		_ = s.eventPublisher.PublishSubscriptionCreated(event)
	}
	return nil
}

func (s *Service) Unsubscribe(sub *subscriptions.Info) error {
	if err := s.repo.UnsubscribeByToken(sub.Token); err != nil {
		return err
	}
	if s.eventPublisher != nil {
		event := events.SubscriptionCancelledEvent{
			SubscriptionID: sub.ID,
			Token:          sub.Token,
			CancelledAt:    time.Now(),
		}
		_ = s.eventPublisher.PublishSubscriptionCancelled(event)
	}
	return nil
}

func (s *Service) Confirm(sub *subscriptions.Info) error {
	if err := s.repo.ConfirmByToken(sub.Token); err != nil {
		return err
	}
	if s.eventPublisher != nil {
		event := events.SubscriptionConfirmedEvent{
			SubscriptionID: sub.ID,
			Token:          sub.Token,
			ConfirmedAt:    time.Now(),
		}
		_ = s.eventPublisher.PublishSubscriptionConfirmed(event)
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
