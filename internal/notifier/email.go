package notifier

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailNotifier struct{}

func (n EmailNotifier) Send(email_to string, message string, subject string) error {
	email_from := os.Getenv("EMAIL_FROM")
	name_from := os.Getenv("EMAIL_FROM_NAME")
	sendgrid_api_key := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail(name_from, email_from)
	email_subject := subject
	to := mail.NewEmail("Recipient", email_to)
	plain_text_content := message
	html_content := message
	m := mail.NewSingleEmail(from, email_subject, to, plain_text_content, html_content)

	client := sendgrid.NewSendClient(sendgrid_api_key)
	response, err := client.Send(m)

	if err != nil {
		fmt.Println("Sending Error:", err)
	} else {
		fmt.Printf("Status Code: %d\n", response.StatusCode)
		fmt.Printf("Response Body: %s\n", response.Body)
		fmt.Printf("Response Headers: %v\n", response.Headers)
	}

	return nil
}
