package config

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	SendGrid    SendGridConfig
	OpenWeather OpenWeatherConfig
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

type OpenWeatherConfig struct {
	ApiKey string
}
