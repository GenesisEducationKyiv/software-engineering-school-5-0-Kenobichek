package notification

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"Weather-Forecast-API/internal/notifier"
	"fmt"
)

type NotificationService interface {
	SendMessage(channelType string, channelValue string, message string, subject string) error
}

func NewNotificationService(notifier notifier.EmailNotifier) *NotificationSender {
	return &NotificationSender{
		notifier: notifier,
	}
}

type NotificationSender struct {
	notifier notifier.EmailNotifier
}

func (s *NotificationSender) SendMessage(
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
