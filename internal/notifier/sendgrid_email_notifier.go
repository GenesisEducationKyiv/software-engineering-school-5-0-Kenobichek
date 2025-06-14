package notifier

import (
	"Weather-Forecast-API/internal/external/sendgrid"
	"fmt"
)

type SendGridEmailNotifier struct {
	sendgrid *external.SendgridNotifier
}

func NewEmailNotifier() (*SendGridEmailNotifier, error) {
	sg, err := external.NewSendgridNotifierFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SendGrid notifier: %w", err)
	}
	return &SendGridEmailNotifier{sendgrid: sg}, nil
}

func (n *SendGridEmailNotifier) Send(to, message, subject string) error {
	target := external.NotificationTarget{
		Type:    "email",
		Address: to,
	}
	return n.sendgrid.Send(target, message, subject)
}
