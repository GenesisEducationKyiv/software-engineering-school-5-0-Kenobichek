package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("config: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func MustLoad() (Config, error) {
	cfg, err := Load()
	if err != nil {
		return cfg, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

func validate(cfg *Config) error {
	errors := []string{}

	if cfg.Server.Port == 0 {
		errors = append(errors, "PORT is required")
	}
	if cfg.Database.Host == "" {
		errors = append(errors, "DB_HOST is required")
	}
	if cfg.Database.User == "" {
		errors = append(errors, "DB_USER is required")
	}
	if cfg.Database.Password == "" {
		errors = append(errors, "DB_PASSWORD is required")
	}
	if cfg.Database.Name == "" {
		errors = append(errors, "DB_NAME is required")
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
	if cfg.Kafka.CommandTopic == "" {
		errors = append(errors, "KAFKA_COMMAND_TOPIC is required")
	}
	if cfg.Kafka.EventTopic == "" {
		errors = append(errors, "KAFKA_EVENT_TOPIC is required")
	}
	if len(errors) > 0 {
		return fmt.Errorf("config validation errors:\n- %s", strings.Join(errors, "\n- "))
	}
	return nil
}
