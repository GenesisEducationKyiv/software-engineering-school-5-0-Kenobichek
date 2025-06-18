package weather_provider

import (
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/internal/weather_provider/models"
	"context"
	"fmt"
)

type OpenWeatherProvider struct {
	geocoding      openweather.GeocodingProvider
	openWeatherAPI openweather.OpenWeatherAPI
}

func NewOpenWeatherProvider(
	geocoding openweather.GeocodingProvider,
	openWeatherAPI openweather.OpenWeatherAPI) OpenWeatherProvider {
	return OpenWeatherProvider{
		geocoding:      geocoding,
		openWeatherAPI: openWeatherAPI,
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
