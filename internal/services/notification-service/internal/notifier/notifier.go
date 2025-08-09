package notifier

import (
	"fmt"
	"strings"

	"notification-service/internal/domain"
)

const (
	ConfirmTemplate = "confirm"
	WeatherUpdateTemplate = "weather_update"
	UnsubscribeTemplate = "unsubscribe"
)

type emailNotifierManager interface {
	Send(to, message, subject string) error
}

type templateRepositoryManager interface {
	GetTemplateByName(name string) (*domain.MessageTemplate, error)
}

type Service struct {
	notifier  emailNotifierManager
	templates templateRepositoryManager
}

func NewService(
	notifier emailNotifierManager,
	templates templateRepositoryManager,
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
	tpl, err := s.templates.GetTemplateByName(ConfirmTemplate)
	if err != nil {
		return fmt.Errorf("failed to load template: %v", err)
	}
	message := strings.ReplaceAll(tpl.Message, "{{ .ConfirmToken }}", token)
	subject := tpl.Subject
	return s.notifier.Send(recipient, message, subject)
}

func (s *Service) SendWeatherUpdate(
	channel string,
	recipient string,
	metrics domain.WeatherMetrics,
) error {
	tpl, err := s.templates.GetTemplateByName(WeatherUpdateTemplate)
	if err != nil {
		return fmt.Errorf("failed to load template: %v", err)
	}
	
	replacements := map[string]string{
		"{{ .City }}":        metrics.City,
		"{{ .Description }}": metrics.Description,
		"{{ .Temperature }}": fmt.Sprintf("%.1f", metrics.Temperature),
		"{{ .Humidity }}":    fmt.Sprintf("%.1f", metrics.Humidity),
	}
			
	message := tpl.Message
	subject := tpl.Subject
	
	for placeholder, value := range replacements {
		message = strings.ReplaceAll(message, placeholder, value)
		if placeholder == "{{ .City }}" {
			subject = strings.ReplaceAll(subject, placeholder, value)
		}
	}
	
	return s.notifier.Send(recipient, message, subject)
}

func (s *Service) SendUnsubscribe(
	channel string,
	recipient string,
	city string,
) error {
	tpl, err := s.templates.GetTemplateByName(UnsubscribeTemplate)
	if err != nil {
		return fmt.Errorf("failed to load template: %v", err)
	}
	message := strings.ReplaceAll(tpl.Message, "{{ .City }}", city)
	subject := strings.ReplaceAll(tpl.Subject, "{{ .City }}", city)

	return s.notifier.Send(recipient, message, subject)
}
