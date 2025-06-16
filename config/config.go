package config

import "fmt"

// Load reads configuration from the specified environment file, populates a Config struct, and validates it.
// Returns the loaded Config or an error if parsing, filling, or validation fails.
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
