package subscribe

import "fmt"

func ConvertFrequency(freq string) (int, error) {
	mins, ok := FrequencyToMinutes()[freq]
	if !ok {
		return 0, fmt.Errorf("invalid frequency: %s", freq)
	}

	return mins, nil
}
