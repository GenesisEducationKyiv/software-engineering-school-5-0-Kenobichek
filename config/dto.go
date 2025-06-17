package config

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	SendGrid    SendGridConfig
	OpenWeather OpenWeatherConfig
}

type ServerConfig struct {
	Port int `envconfig:"PORT" required:"true" default:"8080"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" required:"true" default:"5432"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

type SendGridConfig struct {
	APIKey        string `envconfig:"SENDGRID_API_KEY" required:"true"`
	EmailFrom     string `envconfig:"EMAIL_FROM" required:"true"`
	EmailFromName string `envconfig:"EMAIL_FROM_NAME" required:"true"`
}

type OpenWeatherConfig struct {
	APIKey          string `envconfig:"OPENWEATHERMAP_API_KEY" required:"true"`
	GeocodingAPIURL string `envconfig:"GEOCODING_API_URL" required:"true"`
	WeatherAPIURL   string `envconfig:"WEATHER_API_URL" required:"true"`
}
