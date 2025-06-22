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
		fmt.Printf("Context error: %v\n", err)
		return weather.Metrics{}, err
	}
	if strings.TrimSpace(city) == "" {
		fmt.Printf("Empty city provided\n")
		return weather.Metrics{}, fmt.Errorf("city must not be empty")
	}

	fmt.Printf("Getting coordinates for city: %s\n", city)
	coords, err := wp.geocoding.GetCoordinates(ctx, city)
	if err != nil {
		fmt.Printf("Failed to get coordinates for city %s: %v\n", city, err)
		return weather.Metrics{}, fmt.Errorf("failed to get coordinates: %w", err)
	}
	fmt.Printf("Got coordinates for %s: lat=%f, lon=%f\n", city, coords.Lat, coords.Lon)

	fmt.Printf("Getting weather for coordinates lat=%f, lon=%f\n", coords.Lat, coords.Lon)
	openWeatherData, err := wp.openWeatherAPI.GetWeather(ctx, coords)
	if err != nil {
		fmt.Printf("Failed to get weather for coordinates lat=%f, lon=%f: %v\n", coords.Lat, coords.Lon, err)
		return weather.Metrics{}, fmt.Errorf("failed to get weather: %w", err)
	}

	weatherData := weather.Metrics{
		Temperature: openWeatherData.Temperature,
		Humidity:    openWeatherData.Humidity,
		Description: openWeatherData.Description,
		City:        city,
	}
	fmt.Printf("Got weather for %s: temp=%f, humidity=%f, desc=%s\n",
		city, weatherData.Temperature, weatherData.Humidity, weatherData.Description)

	return weatherData, nil
}
