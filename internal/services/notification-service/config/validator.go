package config

import (
	"fmt"
	"strings"
)

func validate(cfg *Config) error {
	errors := []string{}

	if cfg.Server.Port == 0 {
		errors = append(errors, "PORT is required")
	}
	if len(cfg.Kafka.Brokers) == 0 {
		errors = append(errors, "KAFKA_BROKERS is required")
	} else {
		for i, broker := range cfg.Kafka.Brokers {
			if broker == "" {
				errors = append(errors, fmt.Sprintf("KAFKA_BROKERS[%d] cannot be empty", i))
			}
		}
	}
	if cfg.SendGrid.APIKey == "" {
		errors = append(errors, "SENDGRID_API_KEY is required")
	}
	if cfg.SendGrid.SenderEmail == "" {
		errors = append(errors, "SENDER_EMAIL is required")
	}
	if cfg.SendGrid.SenderName == "" {
		errors = append(errors, "SENDER_NAME is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("config validation errors:\n- %s", strings.Join(errors, "\n- "))
	}
	return nil
}
