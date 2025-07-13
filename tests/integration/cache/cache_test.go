package cache_test

import (
	"Weather-Forecast-API/internal/cache"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/weatherprovider/cached"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWeatherProvider implements the weather provider interface for testing
type mockWeatherProvider struct {
	callCount int
}

func (m *mockWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	m.callCount++
	// Return mock weather data
	return weather.Metrics{
		City:        city,
		Temperature: 20.5 + float64(m.callCount), // Different temperature each call to verify caching
		Humidity:    65.0,
		Description: "Partly cloudy",
	}, nil
}

func TestWeatherCachingIntegration(t *testing.T) {
	// Skip if Redis is not available
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0, 10*time.Minute)
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			t.Logf("Failed to close Redis cache: %v", err)
		}
	}()

	// Clear cache before test to ensure clean state
	ctx := context.Background()
	city := "London"
	_ = redisCache.Delete(ctx, city)

	// Create mock weather provider
	mockProvider := &mockWeatherProvider{}

	// Create cached provider
	cachedProvider := cached.NewCachedWeatherProvider(mockProvider, redisCache)

	// First request - should hit the provider and cache the result
	metrics1, err := cachedProvider.GetWeatherByCity(ctx, city)
	require.NoError(t, err)
	require.NotEmpty(t, metrics1.City)
	require.NotZero(t, metrics1.Temperature)

	// Verify the provider was called once
	assert.Equal(t, 1, mockProvider.callCount)

	// Wait for cache goroutine to finish
	time.Sleep(50 * time.Millisecond)

	// Second request - should hit the cache
	metrics2, err := cachedProvider.GetWeatherByCity(ctx, city)
	require.NoError(t, err)

	// Verify the data is the same (cached result)
	assert.Equal(t, metrics1, metrics2)

	// Verify the provider was not called again (cache hit)
	assert.Equal(t, 1, mockProvider.callCount)

	// Verify cache hit by checking Redis directly
	cachedMetrics, err := redisCache.Get(ctx, city)
	require.NoError(t, err)
	require.NotNil(t, cachedMetrics)
	assert.Equal(t, metrics1, *cachedMetrics)
}
