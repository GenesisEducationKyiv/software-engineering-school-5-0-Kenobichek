package openweatherprovider

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"fmt"
	"strings"
)

type geocodingManager interface {
	GetCoordinates(ctx context.Context, city string) (weather.Coordinates, error)
}

type weatherManager interface {
	GetWeather(ctx context.Context, coords weather.Coordinates) (weather.Metrics, error)
}

type OpenWeatherProvider struct {
	geocoding      geocodingManager
	openWeatherAPI weatherManager
}

func NewOpenWeatherProvider(
	geocoding geocodingManager,
	openWeatherAPI weatherManager,
) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		geocoding:      geocoding,
		openWeatherAPI: openWeatherAPI,
	}
}

func (wp *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	if err := ctx.Err(); err != nil {
		return weather.Metrics{}, err
	}
	if strings.TrimSpace(city) == "" {
		return weather.Metrics{}, fmt.Errorf("city must not be empty")
	}

	coords, err := wp.geocoding.GetCoordinates(ctx, city)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to get coordinates: %w", err)
	}

	openWeatherData, err := wp.openWeatherAPI.GetWeather(ctx, coords)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to get weather: %w", err)
	}

	weatherData := weather.Metrics{
		Temperature: openWeatherData.Temperature,
		Humidity:    openWeatherData.Humidity,
		Description: openWeatherData.Description,
		City:        city,
	}

	return weatherData, nil
}
