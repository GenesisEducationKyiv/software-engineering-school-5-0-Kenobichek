package config

import "fmt"

// Config structures for Subscription Service

type ServerConfig struct {
	Port int `envconfig:"PORT" required:"true" default:"8083"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" required:"true" default:"5432"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

type KafkaConfig struct {
	Brokers    []string `envconfig:"KAFKA_BROKERS" required:"true" default:"kafka:9092"`
	EventTopic string   `envconfig:"KAFKA_EVENT_TOPIC" required:"true" default:"events.subscription"`
	CommandTopic string `envconfig:"KAFKA_COMMAND_TOPIC" required:"true" default:"commands.subscription"`
}

type ObservabilityConfig struct {
	VictoriaMetricsPort int `envconfig:"VICTORIA_METRICS_PORT" required:"true" default:"8428"`
	GrafanaPort int `envconfig:"GRAFANA_PORT" required:"true" default:"3000"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	WeatherServiceAddr string `envconfig:"WEATHER_SERVICE_ADDR" required:"true" default:"weather-service:8081"`
	Observability ObservabilityConfig
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
	)
}
