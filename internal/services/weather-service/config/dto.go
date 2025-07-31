package config

import (
	"time"
)

type Config struct {
	Server      ServerConfig
	OpenWeather OpenWeatherConfig
	WeatherAPI  WeatherAPIConfig
	Redis       RedisConfig
	Monitoring  MonitoringConfig
	Health      HealthConfig
}

type ServerConfig struct {
	Port                    int           `envconfig:"PORT" required:"true" default:"9091"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT" default:"30s"`
}

type OpenWeatherConfig struct {
	APIKey          string `envconfig:"OPENWEATHERMAP_API_KEY" required:"true"`
	GeocodingAPIURL string `envconfig:"GEOCODING_API_URL" required:"true"`
	WeatherAPIURL   string `envconfig:"OPENWEATHERMAP_API_URL" required:"true"`
}

type WeatherAPIConfig struct {
	APIKey string `envconfig:"WEATHER_API_KEY" required:"true"`
	URL    string `envconfig:"WEATHER_API_URL" required:"true" default:"http://api.weatherapi.com/v1/current.json"`
}

type RedisConfig struct {
	Host     string        `envconfig:"REDIS_HOST" required:"true" default:"localhost"`
	Port     int           `envconfig:"REDIS_PORT" required:"true" default:"6379"`
	Password string        `envconfig:"REDIS_PASSWORD"`
	DB       int           `envconfig:"REDIS_DB" default:"0"`
	CacheTTL time.Duration `envconfig:"REDIS_CACHE_TTL" default:"10m"`
}

type MonitoringConfig struct {
	RedisExporterPort       int    `envconfig:"REDIS_EXPORTER_PORT" default:"9121"`
	PrometheusPort          int    `envconfig:"PROMETHEUS_PORT" default:"9090"`
	GrafanaPort             int    `envconfig:"GRAFANA_PORT" default:"3000"`
	GrafanaAdminPassword    string `envconfig:"GRAFANA_ADMIN_PASSWORD" default:"admin"`
	PrometheusRetentionTime string `envconfig:"PROMETHEUS_RETENTION_TIME" default:"200h"`
}

type HealthConfig struct {
	Interval    time.Duration `envconfig:"HEALTH_CHECK_INTERVAL" default:"10s"`
	Retries     int           `envconfig:"HEALTH_CHECK_RETRIES" default:"5"`
	StartPeriod time.Duration `envconfig:"HEALTH_CHECK_START_PERIOD" default:"30s"`
	Timeout     time.Duration `envconfig:"HEALTH_CHECK_TIMEOUT" default:"10s"`
}
