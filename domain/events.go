package domain

type NotificationSentEvent struct {
	NotificationID string
	ChannelType    string
	Recipient      string
	Status         string
	SentAt         int64 // UNIX timestamp
}

type EventPublisherManager interface {
	PublishNotificationSent(event NotificationSentEvent) error
}
