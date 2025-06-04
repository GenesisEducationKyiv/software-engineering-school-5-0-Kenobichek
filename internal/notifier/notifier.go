package notifier

type Notifier interface {
	Send(email_to string, message string, subject string) error
}
