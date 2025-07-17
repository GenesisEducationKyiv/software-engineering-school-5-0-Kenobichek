package provider

import (
	"context"
	"internal/services/weather-service/domain"
)

type weatherAPIManager interface {
	GetWeather(ctx context.Context, city string) (domain.Metrics, error)
}

type WeatherAPIProvider struct {
	weatherapi weatherAPIManager
}

func NewWeatherAPIProvider(weatherapi weatherAPIManager) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		weatherapi: weatherapi,
	}
}

func (wp *WeatherAPIProvider) GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error) {
	return wp.weatherapi.GetWeather(ctx, city)
}
