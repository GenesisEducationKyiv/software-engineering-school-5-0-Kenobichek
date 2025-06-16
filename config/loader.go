package config

import (
	"fmt"
	"strconv"
)

func fillConfig(config *Config, envVars map[string]string) error {
	var err error

	if portStr, exists := envVars["PORT"]; exists && portStr != "" {
		config.Server.Port, err = strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid PORT value: %w", err)
		}
	}

	config.Database.Host = envVars["DB_HOST"]
	config.Database.User = envVars["DB_USER"]
	config.Database.Password = envVars["DB_PASSWORD"]
	config.Database.Name = envVars["DB_NAME"]

	if dbPortStr, exists := envVars["DB_PORT"]; exists && dbPortStr != "" {
		config.Database.Port, err = strconv.Atoi(dbPortStr)
		if err != nil {
			return fmt.Errorf("invalid DB_PORT value: %w", err)
		}
	}

	config.SendGrid.APIKey = envVars["SENDGRID_API_KEY"]
	config.SendGrid.EmailFrom = envVars["EMAIL_FROM"]
	config.SendGrid.EmailFromName = envVars["EMAIL_FROM_NAME"]

	config.Weather.OpenWeatherAPIKey = envVars["OPENWEATHERMAP_API_KEY"]

	return nil
}
