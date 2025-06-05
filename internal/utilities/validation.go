package utilities

import "fmt"

func IsValidChannel(channel string) bool {
	_, ok := SupportedChannels()[channel]

	return ok
}

func ConvertFrequency(freq string) (int, error) {
	mins, ok := FrequencyToMinutes()[freq]
	if !ok {
		return 0, fmt.Errorf("invalid frequency: %s", freq)
	}

	return mins, nil
}
