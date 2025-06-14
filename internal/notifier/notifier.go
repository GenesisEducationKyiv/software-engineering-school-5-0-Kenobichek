package notifier

type Notifier interface {
	Send(to string, message string, subject string) error
}
