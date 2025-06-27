package subscriptions

import "time"

type Info struct {
	ID               int
	ChannelType      string
	ChannelValue     string
	City             string
	FrequencyMinutes int
	Confirmed        bool
	Token            string
	NextNotifiedAt   time.Time
	CreatedAt        time.Time
}
