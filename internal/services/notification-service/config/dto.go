package config

import "time"

type Config struct {
	Server   ServerConfig
	Kafka    KafkaConfig
	SendGrid SendGridConfig
}

type ServerConfig struct {
	Port                    int           `envconfig:"PORT" required:"true" default:"8082"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT" default:"30s"`
}

type KafkaConfig struct {
	Brokers []string `envconfig:"KAFKA_BROKERS" required:"true" default:"kafka:9092"`
}

type SendGridConfig struct {
	APIKey      string `envconfig:"SENDGRID_API_KEY" required:"true"`
	SenderEmail string `envconfig:"SENDER_EMAIL" required:"true"`
	SenderName  string `envconfig:"SENDER_NAME" required:"true"`
}
