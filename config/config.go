package config

import (
	"fmt"
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
func Usage() {
	var cfg Config
	err := envconfig.Usage("", &cfg)
	if err != nil {
		return
	}
}
