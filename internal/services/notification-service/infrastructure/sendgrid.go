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

	fmt.Printf("[SendgridNotifier] Sending email to: %s, subject: %s\n", to, subject)

	resp, err := s.client.Send(m)
	if err != nil {
		fmt.Printf("[SendgridNotifier] Error sending email to %s: %v\n", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("[SendgridNotifier] Sendgrid returned error status code %d for %s. Response body: %s\n", resp.StatusCode, to, resp.Body)
		return fmt.Errorf("sendgrid returned error status code %d", resp.StatusCode)
	}

	fmt.Printf("[SendgridNotifier] Email sent successfully to: %s\n", to)
	return nil
}
