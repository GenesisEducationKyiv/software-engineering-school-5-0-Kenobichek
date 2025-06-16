package sendgrid

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"

	"github.com/sendgrid/sendgrid-go"
)

type NotificationTarget struct {
	Type    string
	Address string
}

type SendgridConfig struct {
	APIKey      string
	SenderEmail string
	SenderName  string
}

func (c *SendgridConfig) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("sendgrid API key is required")
	}
	if c.SenderEmail == "" {
		return fmt.Errorf("sender email is required")
	}
	return nil
}

type SendgridNotifier struct {
	config *SendgridConfig
	client *sendgrid.Client
}

func NewSendgridNotifier(config *SendgridConfig) (*SendgridNotifier, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &SendgridNotifier{
		config: config,
		client: sendgrid.NewSendClient(config.APIKey),
	}, nil
}

func (s *SendgridNotifier) Send(target NotificationTarget, message, subject string) error {
	if target.Type != "email" {
		return fmt.Errorf("invalid notification target for email")
	}
	from := mail.NewEmail(s.config.SenderName, s.config.SenderEmail)
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

func NewSendgridNotifierFromEnv() (*SendgridNotifier, error) {
	config := &SendgridConfig{
		APIKey:      os.Getenv("SENDGRID_API_KEY"),
		SenderEmail: os.Getenv("EMAIL_FROM"),
		SenderName:  os.Getenv("EMAIL_FROM_NAME"),
	}
	return NewSendgridNotifier(config)
}
