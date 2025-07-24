package notifier

import (
	"fmt"
	"strings"

	"notification-service/domain"
)

type EmailNotifierManager interface {
	Send(to, message, subject string) error
}

type TemplateRepositoryManager interface {
	GetTemplateByName(name string) (*domain.MessageTemplate, error)
}

type EventPublisherManager interface {
	PublishNotificationSent(event domain.NotificationSentEvent) error
}

type Service struct {
	notifier  EmailNotifierManager
	templates TemplateRepositoryManager
}

func NewService(
	notifier EmailNotifierManager,
	templates TemplateRepositoryManager,
) *Service {
	return &Service{
		notifier:  notifier,
		templates: templates,
	}
}

func (s *Service) SendConfirmation(
	channel string,
	recipient string,
	token string,
) error {
	tpl, err := s.templates.GetTemplateByName("confirm")
	if err != nil {
		return fmt.Errorf("failed to load template: %v", err)
	}
	message := strings.ReplaceAll(tpl.Message, "{{ confirm_token }}", token)
	subject := tpl.Subject
	return s.notifier.Send(recipient, message, subject)
}

func (s *Service) SendWeatherUpdate(
	channel string,
	recipient string,
	metrics domain.WeatherMetrics,
) error {
	tpl, err := s.templates.GetTemplateByName("weather_update")
	if err != nil {
		return fmt.Errorf("failed to load template: %v", err)
	}
	message := strings.ReplaceAll(tpl.Message, "{{ city }}", metrics.City)
	message = strings.ReplaceAll(message, "{{ description }}", metrics.Description)
	message = strings.ReplaceAll(message, "{{ temperature }}", fmt.Sprintf("%.1f", metrics.Temperature))
	message = strings.ReplaceAll(message, "{{ humidity }}", fmt.Sprintf("%.1f", metrics.Humidity))
	subject := strings.ReplaceAll(tpl.Subject, "{{ city }}", metrics.City)
	return s.notifier.Send(recipient, message, subject)
}

func (s *Service) SendUnsubscribe(
	channel string,
	recipient string,
	city string,
) error {
	tpl, err := s.templates.GetTemplateByName("unsubscribe")
	if err != nil {
		return fmt.Errorf("failed to load template: %v", err)
	}
	message := strings.ReplaceAll(tpl.Message, "{{ city }}", city)
	subject := strings.ReplaceAll(tpl.Subject, "{{ city }}", city)

	return s.notifier.Send(recipient, message, subject)
}
