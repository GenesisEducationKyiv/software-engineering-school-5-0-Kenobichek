package provider

import (
	"context"
	"internal/services/weather-service/domain"
)

type geocodingManager interface {
	GetCoordinates(ctx context.Context, city string) (domain.Coordinates, error)
}

type weatherManager interface {
	GetWeather(ctx context.Context, coords domain.Coordinates) (domain.Metrics, error)
}

type OpenWeatherProvider struct {
	geocoding      geocodingManager
	openWeatherAPI weatherManager
}

func NewOpenWeatherProvider(geocoding geocodingManager, openWeatherAPI weatherManager) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		geocoding:      geocoding,
		openWeatherAPI: openWeatherAPI,
	}
}

func (wp *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error) {
	coords, err := wp.geocoding.GetCoordinates(ctx, city)
	if err != nil {
		return domain.Metrics{}, err
	}
	return wp.openWeatherAPI.GetWeather(ctx, coords)
}
