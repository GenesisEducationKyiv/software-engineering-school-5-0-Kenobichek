package notification

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"Weather-Forecast-API/internal/events"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/emailtemplates"
	"Weather-Forecast-API/internal/templates"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type emailNotifierManager interface {
	Send(to, message, subject string) error
}

type templateRepositoryManager interface {
	GetTemplateByName(name templates.Name) (*emailtemplates.MessageTemplate, error)
}

type eventPublisherManagerher interface {
	PublishNotificationSent(event events.NotificationSentEvent) error
}

func NewService(
	notifier		emailNotifierManager,
	templates		templateRepositoryManager,
	eventPublisher	eventPublisherManagerher,
) *Service {
	return &Service{
		notifier:       notifier,
		templates:      templates,
		eventPublisher: eventPublisher,
	}
}

type Service struct {
	notifier       emailNotifierManager
	templates      templateRepositoryManager
	eventPublisher eventPublisherManagerher
}

func (s *Service) SendConfirmation(
	channel string,
	recipient string,
	token string,
) error {
	switch channel {
	case string(sendgridemailapi.ChannelEmail):
		tpl, err := s.templates.GetTemplateByName(templates.Confirm)

		if err != nil {
			return fmt.Errorf("failed to load template: %v", err)
		}

		message := strings.ReplaceAll(tpl.Message, "{{ confirm_token }}", token)
		subject := tpl.Subject

		return s.notifier.Send(recipient, message, subject)
	default:
		return fmt.Errorf("unsupported channel type: %s", channel)
	}
}

func (s *Service) SendWeatherUpdate(
	channel string,
	recipient string,
	metrics weather.Metrics,
) error {
	switch channel {
	case string(sendgridemailapi.ChannelEmail):
		tpl, err := s.templates.GetTemplateByName(templates.WeatherUpdate)

		if err != nil {
			return fmt.Errorf("failed to load template: %v", err)
		}

		message := strings.ReplaceAll(tpl.Message, "{{ city }}", metrics.City)
		message = strings.ReplaceAll(message, "{{ description }}", metrics.Description)
		message = strings.ReplaceAll(message, "{{ temperature }}", fmt.Sprintf("%.1f", metrics.Temperature))
		message = strings.ReplaceAll(message, "{{ humidity }}", strconv.Itoa(int(metrics.Humidity)))

		subject := strings.ReplaceAll(tpl.Subject, "{{ city }}", metrics.City)

		err = s.notifier.Send(recipient, message, subject)
		if err == nil && s.eventPublisher != nil {
			event := events.NotificationSentEvent{
				NotificationID: "", // You can generate or pass a real ID if available
				ChannelType:    channel,
				Recipient:      recipient,
				Status:         "sent",
				SentAt:         time.Now(),
			}
			_ = s.eventPublisher.PublishNotificationSent(event)
		}
		return err
	default:
		return fmt.Errorf("unsupported channel type: %s", channel)
	}
}

func (s *Service) SendUnsubscribe(
	channel string,
	recipient string,
	city string,
) error {
	switch channel {
	case string(sendgridemailapi.ChannelEmail):
		tpl, err := s.templates.GetTemplateByName(templates.Unsubscribe)

		if err != nil {
			return fmt.Errorf("failed to load template: %v", err)
		}

		message := strings.ReplaceAll(tpl.Message, "{{ city }}", city)
		subject := tpl.Subject

		return s.notifier.Send(recipient, message, subject)
	default:
		return fmt.Errorf("unsupported channel type: %s", channel)
	}
}
