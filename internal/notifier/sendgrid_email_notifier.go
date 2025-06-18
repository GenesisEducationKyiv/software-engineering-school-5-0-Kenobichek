package notifier

import (
	"Weather-Forecast-API/external/sendgrid_email_api"
)

type EmailNotifier interface {
	Send(to, message, subject string) error
}

type SendGridEmailNotifier struct {
	notifier sendgrid_email_api.Notifier
}

func NewSendGridEmailNotifier(notifier sendgrid_email_api.Notifier) SendGridEmailNotifier {
	return SendGridEmailNotifier{
		notifier: notifier,
	}
}

func (n *SendGridEmailNotifier) Send(to, message, subject string) error {
	target := sendgrid_email_api.NotificationTarget{
		Type:    "email",
		Address: to,
	}
	return n.notifier.Send(target, message, subject)
}
