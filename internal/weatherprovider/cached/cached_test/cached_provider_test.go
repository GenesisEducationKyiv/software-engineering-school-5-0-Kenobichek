package cached_test

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/weatherprovider/cached"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWeatherProvider struct {
	GetWeatherByCityFn func(ctx context.Context, city string) (weather.Metrics, error)
}

func (m *mockWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	return m.GetWeatherByCityFn(ctx, city)
}

type mockCache struct {
	GetFn func(ctx context.Context, city string) (*weather.Metrics, error)
	SetFn func(ctx context.Context, city string, metrics weather.Metrics) error
}

func (m *mockCache) Get(ctx context.Context, city string) (*weather.Metrics, error) {
	return m.GetFn(ctx, city)
}

func (m *mockCache) Set(ctx context.Context, city string, metrics weather.Metrics) error {
	return m.SetFn(ctx, city, metrics)
}

func (m *mockCache) Delete(ctx context.Context, city string) error {
	return nil
}

func (m *mockCache) Close() error {
	return nil
}

func TestCachedWeatherProvider_CacheHit(t *testing.T) {
	city := "London"
	cachedMetrics := &weather.Metrics{
		Temperature: 20.5,
		Humidity:    65.0,
		Description: "Partly cloudy",
		City:        city,
	}

	mockProvider := &mockWeatherProvider{
		GetWeatherByCityFn: func(ctx context.Context, city string) (weather.Metrics, error) {
			t.Fatal("Provider should not be called on cache hit")
			return weather.Metrics{}, nil
		},
	}

	mockCache := &mockCache{
		GetFn: func(ctx context.Context, city string) (*weather.Metrics, error) {
			return cachedMetrics, nil
		},
		SetFn: func(ctx context.Context, city string, metrics weather.Metrics) error {
			return nil
		},
	}

	provider := cached.NewCachedWeatherProvider(mockProvider, mockCache)
	ctx := context.Background()

	result, err := provider.GetWeatherByCity(ctx, city)
	require.NoError(t, err)
	assert.Equal(t, *cachedMetrics, result)
}

func TestCachedWeatherProvider_CacheMiss(t *testing.T) {
	city := "Paris"
	expectedMetrics := weather.Metrics{
		Temperature: 18.0,
		Humidity:    70.0,
		Description: "Sunny",
		City:        city,
	}

	providerCalled := false
	mockProvider := &mockWeatherProvider{
		GetWeatherByCityFn: func(ctx context.Context, city string) (weather.Metrics, error) {
			providerCalled = true
			return expectedMetrics, nil
		},
	}

	setCalled := false
	mockCache := &mockCache{
		GetFn: func(ctx context.Context, city string) (*weather.Metrics, error) {
			return nil, nil // Cache miss
		},
		SetFn: func(ctx context.Context, city string, metrics weather.Metrics) error {
			setCalled = true
			assert.Equal(t, expectedMetrics, metrics)
			return nil
		},
	}

	provider := cached.NewCachedWeatherProvider(mockProvider, mockCache)
	ctx := context.Background()

	result, err := provider.GetWeatherByCity(ctx, city)
	require.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
	assert.True(t, providerCalled, "Provider should be called on cache miss")

	// Wait a bit for the async Set operation
	time.Sleep(10 * time.Millisecond)
	assert.True(t, setCalled, "Cache Set should be called")
}

func TestCachedWeatherProvider_CacheError(t *testing.T) {
	city := "Berlin"
	expectedMetrics := weather.Metrics{
		Temperature: 15.0,
		Humidity:    75.0,
		Description: "Cloudy",
		City:        city,
	}

	providerCalled := false
	mockProvider := &mockWeatherProvider{
		GetWeatherByCityFn: func(ctx context.Context, city string) (weather.Metrics, error) {
			providerCalled = true
			return expectedMetrics, nil
		},
	}

	mockCache := &mockCache{
		GetFn: func(ctx context.Context, city string) (*weather.Metrics, error) {
			return nil, errors.New("cache error")
		},
		SetFn: func(ctx context.Context, city string, metrics weather.Metrics) error {
			return nil
		},
	}

	provider := cached.NewCachedWeatherProvider(mockProvider, mockCache)
	ctx := context.Background()

	result, err := provider.GetWeatherByCity(ctx, city)
	require.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
	assert.True(t, providerCalled, "Provider should be called when cache errors")
}

func TestCachedWeatherProvider_ProviderError(t *testing.T) {
	city := "InvalidCity"
	expectedError := errors.New("provider error")

	mockProvider := &mockWeatherProvider{
		GetWeatherByCityFn: func(ctx context.Context, city string) (weather.Metrics, error) {
			return weather.Metrics{}, expectedError
		},
	}

	mockCache := &mockCache{
		GetFn: func(ctx context.Context, city string) (*weather.Metrics, error) {
			return nil, nil // Cache miss
		},
		SetFn: func(ctx context.Context, city string, metrics weather.Metrics) error {
			return nil
		},
	}

	provider := cached.NewCachedWeatherProvider(mockProvider, mockCache)
	ctx := context.Background()

	result, err := provider.GetWeatherByCity(ctx, city)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
	assert.Equal(t, weather.Metrics{}, result)
}
