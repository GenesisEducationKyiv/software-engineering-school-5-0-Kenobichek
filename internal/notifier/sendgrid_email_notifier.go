package notifier

import (
	"Weather-Forecast-API/external/sendgrid"
	"fmt"
)

type SendGridEmailNotifier struct {
	sendgrid *sendgrid.SendgridNotifier
}

func NewEmailNotifier() (*SendGridEmailNotifier, error) {
	sg, err := sendgrid.NewSendgridNotifierFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SendGrid notifier: %w", err)
	}
	return &SendGridEmailNotifier{sendgrid: sg}, nil
}

func (n *SendGridEmailNotifier) Send(to, message, subject string) error {
	target := sendgrid.NotificationTarget{
		Type:    "email",
		Address: to,
	}
	return n.sendgrid.Send(target, message, subject)
}
