package config

import "time"

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	SendGrid    SendGridConfig
	OpenWeather OpenWeatherConfig
	Weather     WeatherConfig
	Redis       RedisConfig
}

type ServerConfig struct {
	Port                    int           `envconfig:"PORT" required:"true" default:"8080"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT" default:"30s"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" required:"true" default:"5432"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

type SendGridConfig struct {
	APIKey      string `envconfig:"SENDGRID_API_KEY" required:"true"`
	SenderEmail string `envconfig:"SENDER_EMAIL" required:"true"`
	SenderName  string `envconfig:"SENDER_NAME" required:"true"`
}

type OpenWeatherConfig struct {
	APIKey          string `envconfig:"OPENWEATHERMAP_API_KEY" required:"true"`
	GeocodingAPIURL string `envconfig:"GEOCODING_API_URL" required:"true"`
	WeatherAPIURL   string `envconfig:"OPENWEATHERMAP_API_URL" required:"true"`
}

type WeatherConfig struct {
	APIKey        string `envconfig:"WEATHER_API_KEY" required:"true"`
	WeatherAPIURL string `envconfig:"WEATHER_API_URL" required:"true"`
}

type RedisConfig struct {
	Host     string        `envconfig:"REDIS_HOST" required:"true" default:"localhost"`
	Port     int           `envconfig:"REDIS_PORT" required:"true" default:"6379"`
	Password string        `envconfig:"REDIS_PASSWORD"`
	DB       int           `envconfig:"REDIS_DB" default:"0"`
	TTL      time.Duration `envconfig:"REDIS_CACHE_TTL" default:"10m"`
}
