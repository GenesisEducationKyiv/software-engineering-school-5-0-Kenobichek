package chain_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/weatherprovider/chain"
)

type mockProvider struct {
	metrics weather.Metrics
	err     error
}

func (m *mockProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	return m.metrics, m.err
}

func TestChainWeatherProvider_Success(t *testing.T) {
	mock := &mockProvider{metrics: weather.Metrics{Temperature: 25.0}, err: nil}
	cwp := chain.NewChainWeatherProvider(mock)

	metrics, err := cwp.GetWeatherByCity(context.Background(), "London")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if metrics.Temperature != 25.0 {
		t.Errorf("expected temperature 25.0, got %v", metrics.Temperature)
	}
}

func TestChainWeatherProvider_Fallback(t *testing.T) {
	primary := &mockProvider{err: errors.New("fail")}
	fallback := &mockProvider{metrics: weather.Metrics{Temperature: 15.0}, err: nil}

	cwp := chain.NewChainWeatherProvider(primary)
	cwp.SetNext(chain.NewChainWeatherProvider(fallback))

	metrics, err := cwp.GetWeatherByCity(context.Background(), "Paris")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if metrics.Temperature != 15.0 {
		t.Errorf("expected temperature 15.0, got %v", metrics.Temperature)
	}
}

func TestChainWeatherProvider_NoFallback(t *testing.T) {
	primary := &mockProvider{err: errors.New("fail")}
	cwp := chain.NewChainWeatherProvider(primary)

	_, err := cwp.GetWeatherByCity(context.Background(), "Berlin")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected error to contain original error message, got: %v", err)
	}
}
