package cache_test

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/internal/cache"
	"Weather-Forecast-API/internal/weatherprovider/cached"
	"Weather-Forecast-API/internal/weatherprovider/openweatherprovider"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	// Load config
	cfg, err := config.MustLoad()
	require.NoError(t, err)

	// Create HTTP client
	httpClient := &http.Client{Timeout: 10 * time.Second}

	// Create geocoding service
	geoSvc := openweather.NewGeocodingService(
		httpClient,
		cfg.OpenWeather.GeocodingAPIURL,
		cfg.OpenWeather.APIKey,
	)

	// Create weather API
	owAPI := openweather.NewOpenWeatherAPI(
		httpClient,
		cfg.OpenWeather.WeatherAPIURL,
		cfg.OpenWeather.APIKey,
	)

	// Create weather provider
	weatherProvider := openweatherprovider.NewOpenWeatherProvider(geoSvc, owAPI)

	// Create cached provider
	cachedProvider := cached.NewCachedWeatherProvider(weatherProvider, redisCache)

	ctx := context.Background()
	city := "London"

	// First request - should hit the provider and cache the result
	start := time.Now()
	metrics1, err := cachedProvider.GetWeatherByCity(ctx, city)
	firstRequestTime := time.Since(start)
	require.NoError(t, err)
	require.NotEmpty(t, metrics1.City)
	require.NotZero(t, metrics1.Temperature)

	// Second request - should hit the cache
	start = time.Now()
	metrics2, err := cachedProvider.GetWeatherByCity(ctx, city)
	secondRequestTime := time.Since(start)
	require.NoError(t, err)

	// Verify the data is the same
	assert.Equal(t, metrics1, metrics2)

	// Verify the second request was faster (cache hit)
	assert.Less(t, secondRequestTime, firstRequestTime, "Cached request should be faster")

	// Verify cache hit by checking Redis directly
	cachedMetrics, err := redisCache.Get(ctx, city)
	require.NoError(t, err)
	require.NotNil(t, cachedMetrics)
	assert.Equal(t, metrics1, *cachedMetrics)
}
