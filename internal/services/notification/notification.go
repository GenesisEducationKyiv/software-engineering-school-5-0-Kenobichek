package notification

import (
	"Weather-Forecast-API/external/sendgrid_email_api"
	"Weather-Forecast-API/internal/notifier"
	"fmt"
)

type NotificationService interface {
	SendMessage(channelType string, channelValue string, message string, subject string) error
}

func NewNotificationService(notifier sendgrid_email_api.Notifier) NotificationSender {
	return NotificationSender{
		notifier: notifier,
	}
}

type NotificationSender struct {
	notifier sendgrid_email_api.Notifier
}

func (s *NotificationSender) SendMessage(
	channelType string,
	channelValue string,
	message string,
	subject string) error {

	switch channelType {
	case "email":
		n := notifier.NewSendGridEmailNotifier(s.notifier)
		return n.Send(channelValue, message, subject)
	default:
		return fmt.Errorf("unsupported channel type: %s", channelType)
	}
}
