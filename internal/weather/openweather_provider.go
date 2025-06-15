package weather

import (
	"Weather-Forecast-API/internal/external/openweather"
	"Weather-Forecast-API/internal/weather/models"
	"context"
	"fmt"
)

type OpenWeatherProvider struct {
	geocoding      *openweather.GeocodingService
	openWeatherAPI *openweather.OpenWeatherAPI
}

func NewOpenWeatherProvider(apiKey string) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		geocoding:      openweather.NewGeocodingService(apiKey),
		openWeatherAPI: openweather.NewWeatherService(apiKey),
	}
}

func (wp *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (models.WeatherData, error) {
	coords, err := wp.geocoding.GetCoordinates(ctx, city)
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("failed to get coordinates: %w", err)
	}

	openWeatherData, err := wp.openWeatherAPI.GetWeather(ctx, coords)
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("failed to get weather: %w", err)
	}

	weatherData := models.WeatherData{
		Temperature: openWeatherData.Temperature,
		Humidity:    openWeatherData.Humidity,
		Description: openWeatherData.Description,
	}

	return weatherData, nil
}
