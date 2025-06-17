package notifier

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/external/sendgrid"
	"fmt"
)

type EmailNotifier interface {
	Send(to, message, subject string) error
}

type SendGridEmailNotifier struct {
	sendgrid sendgrid.Notifier
	cfg      *config.Config
}

var newSendgridNotifier = sendgrid.NewSendgridNotifier

func NewSendGridEmailNotifier(cfg *config.Config) (SendGridEmailNotifier, error) {
	sg, err := newSendgridNotifier(cfg)
	if err != nil {
		return SendGridEmailNotifier{}, fmt.Errorf("failed to initialize SendGrid notifier: %w", err)
	}
	return SendGridEmailNotifier{
		sendgrid: sg,
		cfg:      cfg,
	}, nil
}

func (n SendGridEmailNotifier) Send(to, message, subject string) error {
	target := sendgrid.NotificationTarget{
		Type:    "email",
		Address: to,
	}
	return n.sendgrid.Send(target, message, subject)
}
