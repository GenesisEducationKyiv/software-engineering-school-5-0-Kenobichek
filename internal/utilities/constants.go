package utilities

const (
	hourlyMinutes = 60
	dailyMinutes  = 1440
)

func SupportedChannels() map[string]struct{} {
	return map[string]struct{}{
		"email": {},
		// "sms":   {},
	}
}

func FrequencyToMinutes() map[string]int {
	return map[string]int{
		"hourly": hourlyMinutes,
		"daily":  dailyMinutes,
	}
}
