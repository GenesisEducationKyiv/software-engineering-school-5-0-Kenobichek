package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SendGrid SendGridConfig
	Weather  WeatherConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type SendGridConfig struct {
	APIKey        string
	EmailFrom     string
	EmailFromName string
}

type WeatherConfig struct {
	OpenWeatherAPIKey string
}
