package config

import "fmt"

func Load(envPath string) (*Config, error) {
	envVars, err := parseEnvFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing .env file: %w", err)
	}

	config := &Config{}

	if err := fillConfig(config, envVars); err != nil {
		return nil, fmt.Errorf("error filling config: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return config, nil
}
