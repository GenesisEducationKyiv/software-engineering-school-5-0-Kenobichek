package notifier

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"fmt"
)

type EmailNotifier interface {
	Send(to, message, subject string) error
}

type SendGridEmailNotifier struct {
	notifier sendgridemailapi.Notifier
}

func NewSendGridEmailNotifier(notifier sendgridemailapi.Notifier) *SendGridEmailNotifier {
	return &SendGridEmailNotifier{
		notifier: notifier,
	}
}

func (n *SendGridEmailNotifier) Send(to, message, subject string) error {
	target := sendgridemailapi.NotificationTarget{
		Type:    sendgridemailapi.ChannelEmail,
		Address: to,
	}
	if err := n.notifier.Send(target, message, subject); err != nil {
		return fmt.Errorf("failed to send email notification to %s: %w", to, err)
	}
	return nil
}
