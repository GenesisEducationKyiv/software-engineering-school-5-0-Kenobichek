package utilities

var SupportedChannels = map[string]struct{}{
	"email": {},
	// "sms":   {},
}

var FrequencyToMinutes = map[string]int{
	"hourly": 60,
	"daily":  1440,
}
