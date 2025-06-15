package openweather

import (
	"context"
	"fmt"
)

type OpenWeatherProvider struct {
	geocoding      *GeocodingService
	openWeatherAPI *OpenWeatherAPI
}

func NewOpenWeatherProvider(apiKey string) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		geocoding:      NewGeocodingService(apiKey),
		openWeatherAPI: NewWeatherService(apiKey),
	}
}

func (wp *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (WeatherData, error) {
	coords, err := wp.geocoding.GetCoordinates(ctx, city)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to get coordinates: %w", err)
	}

	weatherData, err := wp.openWeatherAPI.GetWeather(ctx, coords)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to get weather: %w", err)
	}

	return weatherData, nil
}
