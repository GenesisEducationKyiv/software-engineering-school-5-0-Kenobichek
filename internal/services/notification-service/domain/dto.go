package domain

type MessageTemplate struct {
	Subject string
	Message string
}

type NotificationSentEvent struct {
	NotificationID string
	ChannelType    string
	Recipient      string
	Status         string
	SentAt         int64
}

type WeatherMetrics struct {
	City        string
	Description string
	Temperature float64
	Humidity    int
}
