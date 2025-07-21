package handlers

import (
	"fmt"
	"strings"
)

func validateSubscriptionParams(email, city, frequency string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if city == "" {
		return fmt.Errorf("city is required")
	}
	if frequency != "hourly" && frequency != "daily" {
		return fmt.Errorf("invalid frequency: must be 'hourly' or 'daily'")
	}
	return nil
}

func frequencyToMinutes(frequency string) (int, error) {
	switch frequency {
	case "hourly":
		return 60, nil
	case "daily":
		return 1440, nil
	default:
		return 0, fmt.Errorf("invalid frequency")
	}
}

func validateConfirmSubscriptionParams(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

func validateUnsubscribeParams(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

func validateWeatherParams(city string) error {
	if city == "" {
		return fmt.Errorf("city parameter is required")
	}
	return nil
}
