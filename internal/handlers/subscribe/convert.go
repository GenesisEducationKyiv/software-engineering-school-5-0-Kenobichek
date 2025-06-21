package subscribe

import "fmt"

func convertFrequencyToMinutes(freq string) (int, error) {
	mins, ok := frequencyToMinutes()[freq]
	if !ok {
		return 0, fmt.Errorf("invalid frequency: %s", freq)
	}

	return mins, nil
}
