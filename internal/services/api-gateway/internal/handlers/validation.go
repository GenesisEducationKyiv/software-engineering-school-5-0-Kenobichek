package handlers

import (
	"fmt"
	"regexp"
)

const (
	hourlyFrequencyMinutes = 60
	dailyFrequencyMinutes = 1440
	maxCityNameLength = 100
	maxEmailLength = 100
	regexEmail = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func validateSubscriptionParams(email, city, frequency string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if matched, _ := regexp.MatchString(regexEmail, email); !matched {
		return fmt.Errorf("invalid email format")
	}
	if len(email) > maxEmailLength {
		return fmt.Errorf("email too long")
	}
	if len(city) > maxCityNameLength {
		return fmt.Errorf("city name too long")
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
		return hourlyFrequencyMinutes, nil
	case "daily":
		return dailyFrequencyMinutes, nil
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
