package sendgrid

import (
	"Weather-Forecast-API/config"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Notifier interface {
	Send(target NotificationTarget, message, subject string) error
}

type SendgridNotifier struct {
	cfg    *config.Config
	client *sendgrid.Client
}

func NewSendgridNotifier(cfg *config.Config) (*SendgridNotifier, error) {
	return &SendgridNotifier{
		cfg:    cfg,
		client: sendgrid.NewSendClient(cfg.SendGrid.APIKey),
	}, nil
}

func (s *SendgridNotifier) Send(target NotificationTarget, message, subject string) error {
	if target.Type != "email" {
		return fmt.Errorf("invalid notification target for email")
	}
	from := mail.NewEmail(s.cfg.SendGrid.EmailFromName, s.cfg.SendGrid.EmailFrom)
	to := mail.NewEmail("", target.Address)
	m := mail.NewSingleEmail(from, subject, to, message, message)

	resp, err := s.client.Send(m)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("sendgrid returned status code %d", resp.StatusCode)
	}
	return nil
}
