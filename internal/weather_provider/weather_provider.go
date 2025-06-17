package weather_provider

import (
	"Weather-Forecast-API/internal/weather_provider/models"
	"context"
)

type WeatherProvider interface {
	GetWeatherByCity(ctx context.Context, city string) (models.WeatherData, error)
}
