package sendgrid

type NotificationTarget struct {
	Type    string
	Address string
}

type Config struct {
	APIKey      string
	SenderEmail string
	SenderName  string
}
