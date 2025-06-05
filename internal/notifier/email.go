package notifier

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailNotifier struct{}

func (n EmailNotifier) Send(emailTo string, message string, subject string) error {
	emailFrom := os.Getenv("EMAIL_FROM")
	nameFrom := os.Getenv("EMAIL_FROM_NAME")
	apiKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail(nameFrom, emailFrom)
	emailSubject := subject
	to := mail.NewEmail("Recipient", emailTo)
	plainTextContent := message
	htmlContent := message
	m := mail.NewSingleEmail(from, emailSubject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(m)

	if err != nil {
		return err
	}

	return nil
}
