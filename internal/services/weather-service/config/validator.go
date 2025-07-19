package config

import "fmt"

func validate(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	if cfg.Server.GracefulShutdownTimeout < 0 {
		return fmt.Errorf("graceful shutdown timeout must be >= 0")
	}
	if cfg.OpenWeather.APIKey == "" {
		return fmt.Errorf("OPENWEATHERMAP_API_KEY is required")
	}
	if cfg.OpenWeather.GeocodingAPIURL == "" {
		return fmt.Errorf("GEOCODING_API_URL is required")
	}
	if cfg.OpenWeather.WeatherAPIURL == "" {
		return fmt.Errorf("OPENWEATHERMAP_API_URL is required")
	}
	if cfg.WeatherAPI.APIKey == "" {
		return fmt.Errorf("WEATHER_API_KEY is required")
	}
	if cfg.WeatherAPI.URL == "" {
		return fmt.Errorf("WEATHER_API_URL is required")
	}
	if cfg.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}
	if cfg.Redis.Port <= 0 || cfg.Redis.Port > 65535 {
		return fmt.Errorf("invalid redis port: %d", cfg.Redis.Port)
	}
	if cfg.Redis.DB < 0 {
		return fmt.Errorf("redis db must be >= 0")
	}
	if cfg.Redis.CacheTTL < 0 {
		return fmt.Errorf("redis cache ttl must be >= 0")
	}
	return nil
}
