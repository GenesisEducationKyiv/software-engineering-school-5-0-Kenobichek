package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               string   `envconfig:"PORT" default:"8084"`
	KafkaBrokers       []string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaTopic         string   `envconfig:"KAFKA_TOPIC" required:"true"`
	WeatherServiceAddr string   `envconfig:"WEATHER_SERVICE_ADDR" default:"weather-service:8083"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}
