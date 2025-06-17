package notification

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/notifier"
	"fmt"
)

type NotificationService interface {
	SendMessage(channelType string, channelValue string, message string, subject string) error
}

func NewNotificationService(cfg *config.Config) NotificationService {
	return &notificationService{cfg: cfg}
}

type notificationService struct {
	cfg *config.Config
}

func (s *notificationService) SendMessage(
	channelType string,
	channelValue string,
	message string,
	subject string) error {

	switch channelType {
	case "email":
		n, err := notifier.NewSendGridEmailNotifier(s.cfg)
		if err != nil {
			return err
		}
		return n.Send(channelValue, message, subject)
	default:
		return fmt.Errorf("unsupported channel type: %s", channelType)
	}
}
