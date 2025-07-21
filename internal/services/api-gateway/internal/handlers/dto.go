package handlers

type SubscriptionCommand struct {
	Command          string `json:"command"`
	ChannelType      string `json:"channel_type,omitempty"`
	ChannelValue     string `json:"channel_value,omitempty"`
	City             string `json:"city,omitempty"`
	Frequency        string `json:"frequency,omitempty"`
	FrequencyMinutes int    `json:"frequency_minutes,omitempty"`
	Token            string `json:"token,omitempty"`
}