package sendgrid_email_api

type NotificationTarget struct {
	Type    string
	Address string
}

type Config struct {
	APIKey      string
	SenderEmail string
	SenderName  string
}
