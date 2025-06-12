package notifier

type Notifier interface {
	Send(emailTo string, message string, subject string) error
}
