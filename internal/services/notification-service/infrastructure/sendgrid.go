package infrastructure

import (
	"fmt"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridNotifier struct {
	client      sendgridClientManager
	senderName  string
	senderEmail string
}

type sendgridClientManager interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

func NewSendgridNotifier(
	client sendgridClientManager,
	senderName string,
	senderEmail string,
) *SendgridNotifier {
	return &SendgridNotifier{
		client:      client,
		senderName:  senderName,
		senderEmail: senderEmail,
	}
}

func (s *SendgridNotifier) Send(to, message, subject string) error {
	from := mail.NewEmail(s.senderName, s.senderEmail)
	toEmail := mail.NewEmail("", to)
	m := mail.NewSingleEmail(from, subject, toEmail, message, message)

	resp, err := s.client.Send(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("sendgrid returned error status code %d", resp.StatusCode)
	}
	return nil
}
