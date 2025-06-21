package sengridnotifier

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"fmt"
)

type notifierManager interface {
	Send(target sendgridemailapi.NotificationTarget, message, subject string) error
}

type SendGridEmailNotifier struct {
	notifier notifierManager
}

func NewSendGridEmailNotifier(notifier notifierManager) *SendGridEmailNotifier {
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
