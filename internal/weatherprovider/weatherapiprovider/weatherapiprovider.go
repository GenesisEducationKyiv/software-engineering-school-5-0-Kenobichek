package weatherapiprovider

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"fmt"
	"strings"
)

type weatherManager interface {
	GetWeather(ctx context.Context, city string) (weather.Metrics, error)
}

type WeatherAPIProvider struct {
	weatherapi weatherManager
}

func NewWeatherAPIProvider(
	weatherapi weatherManager,
) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		weatherapi: weatherapi,
	}
}

func (wp *WeatherAPIProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	if err := ctx.Err(); err != nil {
		return weather.Metrics{}, err
	}
	if strings.TrimSpace(city) == "" {
		return weather.Metrics{}, fmt.Errorf("city must not be empty")
	}

	metrics, err := wp.weatherapi.GetWeather(ctx, city)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to get weather: %w", err)
	}

	weatherdata := weather.Metrics{
		Temperature: metrics.Temperature,
		Humidity:    metrics.Humidity,
		Description: metrics.Description,
		City:        city,
	}

	return weatherdata, nil
}
