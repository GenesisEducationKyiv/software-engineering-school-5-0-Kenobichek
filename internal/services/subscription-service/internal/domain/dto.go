package domain

type SubscriptionCommand struct {
	Command          string `json:"command"` // subscribe, confirm, unsubscribe
	ChannelType      string `json:"channel_type"`
	ChannelValue     string `json:"channel_value"`
	City             string `json:"city"`
	FrequencyMinutes int    `json:"frequency_minutes"`
	Token            string `json:"token"`
}

type SubscriptionEvent struct {
	EventType        string `json:"event_type"`
	ChannelType      string `json:"channel_type"`
	ChannelValue     string `json:"channel_value"`
	City             string `json:"city"`
	FrequencyMinutes int    `json:"frequency_minutes,omitempty"`
	Token            string `json:"token,omitempty"`
}

type WeatherMetrics struct {
	City        string
	Description string
	Temperature float64
	Humidity    float64
}

type WeatherUpdateEvent struct {
	Metrics   WeatherMetrics	`json:"metrics"`
	UpdatedAt int64				`json:"updated_at"`
	Email     string			`json:"channel_value"`
}