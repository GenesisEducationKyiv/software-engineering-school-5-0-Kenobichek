package weather

import (
	"Weather-Forecast-API/internal/weather/models"
	"context"
)

type WeatherProvider interface {
	GetWeatherByCity(ctx context.Context, city string) (models.WeatherData, error)
}
