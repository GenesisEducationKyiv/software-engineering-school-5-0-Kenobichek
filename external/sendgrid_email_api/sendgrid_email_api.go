package sendgrid_email_api

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/constants"
	"fmt"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Notifier interface {
	Send(target NotificationTarget, message, subject string) error
}

type SendgridNotifier struct {
	cfg    *config.Config
	client SendGridClient
}

type SendGridClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

func NewSendgridNotifier(client SendGridClient, cfg *config.Config) SendgridNotifier {
	return SendgridNotifier{
		client: client,
		cfg:    cfg,
	}
}

func (s *SendgridNotifier) Send(target NotificationTarget, message, subject string) error {
	if target.Type != constants.ChannelEmail {
		return fmt.Errorf("invalid notification target type %s, expected email", target.Type)
	}

	from := mail.NewEmail(s.cfg.SendGrid.SenderName, s.cfg.SendGrid.SenderEmail)
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