package notification

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"fmt"
)

type emailNotifierManager interface {
	Send(to, message, subject string) error
}

func NewService(notifier emailNotifierManager) *Service {
	return &Service{
		notifier: notifier,
	}
}

type Service struct {
	notifier emailNotifierManager
}

func (s *Service) SendMessage(
	channelType string,
	channelValue string,
	message string,
	subject string) error {

	switch channelType {
	case string(sendgridemailapi.ChannelEmail):
		return s.notifier.Send(channelValue, message, subject)
	default:
		return fmt.Errorf("unsupported channel type: %s", channelType)
	}
}
