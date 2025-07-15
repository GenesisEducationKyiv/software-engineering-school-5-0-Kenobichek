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
