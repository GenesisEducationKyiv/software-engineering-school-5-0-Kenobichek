package weather_provider

import (
	"Weather-Forecast-API/external/openweather"
	"context"
	"fmt"
	"strings"
)

type OpenWeatherProvider struct {
	geocoding      openweather.GeocodingProvider
	openWeatherAPI openweather.WeatherProvider
}

func NewOpenWeatherProvider(
	geocoding openweather.GeocodingProvider,
	openWeatherAPI openweather.WeatherProvider) OpenWeatherProvider {
	return OpenWeatherProvider{
		geocoding:      geocoding,
		openWeatherAPI: openWeatherAPI,
	}
}

func (wp *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (openweather.WeatherData, error) {
	if err := ctx.Err(); err != nil {
		return openweather.WeatherData{}, err
	}
	if strings.TrimSpace(city) == "" {
		return openweather.WeatherData{}, fmt.Errorf("city must not be empty")
	}

	coords, err := wp.geocoding.GetCoordinates(ctx, city)
	if err != nil {
		return openweather.WeatherData{}, fmt.Errorf("failed to get coordinates: %w", err)
	}

	openWeatherData, err := wp.openWeatherAPI.GetWeather(ctx, coords)
	if err != nil {
		return openweather.WeatherData{}, fmt.Errorf("failed to get weather: %w", err)
	}

	weatherData := openweather.WeatherData{
		Temperature: openWeatherData.Temperature,
		Humidity:    openWeatherData.Humidity,
		Description: openWeatherData.Description,
	}

	return weatherData, nil
}
