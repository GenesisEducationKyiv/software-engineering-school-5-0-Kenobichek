package provider_test

import (
	"context"
	"errors"
	"testing"

	"internal/services/weather-service/internal/domain"
	provider "internal/services/weather-service/internal/provider"
)

type mockGeocodingManager struct {
	coords domain.Coordinates
	err    error
}

func (m *mockGeocodingManager) GetCoordinates(ctx context.Context, city string) (domain.Coordinates, error) {
	return m.coords, m.err
}

type mockWeatherManager struct {
	metrics domain.Metrics
	err     error
}

func (m *mockWeatherManager) GetWeather(ctx context.Context, coords domain.Coordinates) (domain.Metrics, error) {
	return m.metrics, m.err
}

func TestOpenWeatherProvider_GetWeatherByCity_Success(t *testing.T) {
	geo := &mockGeocodingManager{
		coords: domain.Coordinates{Lat: 50.45, Lon: 30.52},
	}
	weather := &mockWeatherManager{
		metrics: domain.Metrics{
			Temperature: 20.5,
			Humidity:    60,
			Description: "Clear",
			City:        "Kyiv",
		},
	}
	providerInstance := provider.NewOpenWeatherProvider(geo, weather)

	result, err := providerInstance.GetWeatherByCity(context.Background(), "Kyiv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Temperature != 20.5 || result.City != "Kyiv" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestOpenWeatherProvider_GetWeatherByCity_GeocodingError(t *testing.T) {
	geo := &mockGeocodingManager{
		err: errors.New("geocoding failed"),
	}
	weather := &mockWeatherManager{}
	providerInstance := provider.NewOpenWeatherProvider(geo, weather)

	_, err := providerInstance.GetWeatherByCity(context.Background(), "UnknownCity")
	if err == nil || err.Error() != "geocoding failed" {
		t.Errorf("expected geocoding error, got: %v", err)
	}
}

func TestOpenWeatherProvider_GetWeatherByCity_WeatherError(t *testing.T) {
	geo := &mockGeocodingManager{
		coords: domain.Coordinates{Lat: 1, Lon: 2},
	}
	weather := &mockWeatherManager{
		err: errors.New("weather fetch failed"),
	}
	providerInstance := provider.NewOpenWeatherProvider(geo, weather)

	_, err := providerInstance.GetWeatherByCity(context.Background(), "AnyCity")
	if err == nil || err.Error() != "weather fetch failed" {
		t.Errorf("expected weather error, got: %v", err)
	}
}
