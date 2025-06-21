package sendgridemailapi

import (
	"fmt"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Notifier interface {
	Send(target NotificationTarget, message, subject string) error
}

type SendgridNotifier struct {
	client      SendGridClient
	senderName  string
	senderEmail string
}

type SendGridClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

func NewSendgridNotifier(
	client SendGridClient,
	senderName string,
	senderEmail string) *SendgridNotifier {
	return &SendgridNotifier{
		client:      client,
		senderName:  senderName,
		senderEmail: senderEmail,
	}
}

func (s *SendgridNotifier) Send(target NotificationTarget, message, subject string) error {
	if target.Type != ChannelEmail {
		return fmt.Errorf("invalid notification target type %s, expected email", target.Type)
	}

	from := mail.NewEmail(s.senderName, s.senderEmail)
	to := mail.NewEmail("", target.Address)
	m := mail.NewSingleEmail(from, subject, to, message, message)

	resp, err := s.client.Send(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("sendgrid returned error status code %d", resp.StatusCode)
	}
	return nil
}
