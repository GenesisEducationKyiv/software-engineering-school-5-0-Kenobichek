package weather_provider

import (
	"Weather-Forecast-API/external/openweather"
	"context"
)

type WeatherProvider interface {
	GetWeatherByCity(ctx context.Context, city string) (openweather.WeatherData, error)
}
