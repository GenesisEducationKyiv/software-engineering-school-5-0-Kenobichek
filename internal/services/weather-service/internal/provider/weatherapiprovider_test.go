package provider_test

import (
	"context"
	"errors"
	"testing"

	"internal/services/weather-service/internal/domain"
	provider "internal/services/weather-service/internal/provider"
)

type mockWeatherAPI struct {
	metrics domain.Metrics
	err     error
}

func (m *mockWeatherAPI) GetWeather(ctx context.Context, city string) (domain.Metrics, error) {
	return m.metrics, m.err
}

func TestWeatherAPIProvider_Success(t *testing.T) {
	weather := &mockWeatherAPI{metrics: domain.Metrics{City: "Kyiv"}}
	prov := provider.NewWeatherAPIProvider(weather)

	result, err := prov.GetWeatherByCity(context.Background(), "Kyiv")
	if err != nil || result.City != "Kyiv" {
		t.Errorf("expected success, got: %+v, err: %v", result, err)
	}
}

func TestWeatherAPIProvider_Error(t *testing.T) {
	weather := &mockWeatherAPI{err: errors.New("api error")}
	prov := provider.NewWeatherAPIProvider(weather)

	_, err := prov.GetWeatherByCity(context.Background(), "Kyiv")
	if err == nil || err.Error() != "api error" {
		t.Errorf("expected api error, got: %v", err)
	}
}
