package service

import (
	"Weather-Forecast-API/internal/notifier"
	"fmt"
)

type NotificationService interface {
	SendMessage(channelType string, channelValue string, message string, subject string) error
}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

type notificationService struct{}

func (s *notificationService) SendMessage(
	channelType string,
	channelValue string,
	message string,
	subject string) error {

	switch channelType {
	case "email":
		n, err := notifier.NewEmailNotifier()
		if err != nil {
			return err
		}
		return n.Send(channelValue, message, subject)
	default:
		return fmt.Errorf("unsupported channel type: %s", channelType)
	}
}
