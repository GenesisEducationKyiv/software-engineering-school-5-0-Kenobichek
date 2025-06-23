package config

import (
	"fmt"
	"strings"
)

func validate(config *Config) error {
	var errors []string

	if config.Server.Port == 0 {
		errors = append(errors, "PORT is required")
	} else if config.Server.Port < 1 || config.Server.Port > 65535 {
		errors = append(errors, "PORT must be between 1 and 65535")
	}

	if config.Database.Host == "" {
		errors = append(errors, "DB_HOST is required")
	}
	if config.Database.Port == 0 {
		errors = append(errors, "DB_PORT is required")
	} else if config.Database.Port < 1 || config.Database.Port > 65535 {
		errors = append(errors, "DB_PORT must be between 1 and 65535")
	}
	if config.Database.User == "" {
		errors = append(errors, "DB_USER is required")
	}
	if config.Database.Password == "" {
		errors = append(errors, "DB_PASSWORD is required")
	}
	if config.Database.Name == "" {
		errors = append(errors, "DB_NAME is required")
	}

	if config.SendGrid.APIKey != "" {
		if config.SendGrid.SenderEmail == "" {
			errors = append(errors, "EMAIL_FROM is required when SENDGRID_API_KEY is set")
		}
		if config.SendGrid.SenderName == "" {
			errors = append(errors, "EMAIL_FROM_NAME is required when SENDGRID_API_KEY is set")
		}
	}

	if config.OpenWeather.APIKey == "" {
		errors = append(errors, "OPENWEATHERMAP_API_KEY is required")
	}
	if config.OpenWeather.WeatherAPIURL == "" {
		errors = append(errors, "WEATHER_API_URL is required")
	}
	if config.OpenWeather.GeocodingAPIURL == "" {
		errors = append(errors, "GEOCODING_API_URL is required")
	}
	if len(errors) > 0 {
		return fmt.Errorf("config validation errors:\n- %s", strings.Join(errors, "\n- "))
	}

	return nil
}
