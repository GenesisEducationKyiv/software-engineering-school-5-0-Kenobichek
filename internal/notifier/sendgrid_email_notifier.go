package notifier

import (
	"Weather-Forecast-API/external/sendgrid_email_api"
	"fmt"
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
		Type:    sendgrid_email_api.ChannelEmail,
		Address: to,
	}
	if err := n.notifier.Send(target, message, subject); err != nil {
		return fmt.Errorf("failed to send email notification to %s: %w", to, err)
	}
	return nil
}
