package cache_test

import (
	"Weather-Forecast-API/internal/cache"
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisCache_GetSet(t *testing.T) {
	// Skip if Redis is not available
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0, 10*time.Minute)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			t.Logf("Failed to close Redis cache: %v", err)
		}
	}()

	ctx := context.Background()
	city := "London"
	expectedMetrics := weather.Metrics{
		Temperature: 20.5,
		Humidity:    65.0,
		Description: "Partly cloudy",
		City:        city,
	}

	// Test Set
	err = redisCache.Set(ctx, city, expectedMetrics)
	require.NoError(t, err)

	// Test Get
	cachedMetrics, err := redisCache.Get(ctx, city)
	require.NoError(t, err)
	require.NotNil(t, cachedMetrics)
	assert.Equal(t, expectedMetrics, *cachedMetrics)
}

func TestRedisCache_GetNonExistent(t *testing.T) {
	// Skip if Redis is not available
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0, 10*time.Minute)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			t.Logf("Failed to close Redis cache: %v", err)
		}
	}()

	ctx := context.Background()
	city := "NonExistentCity"

	// Test Get for non-existent key
	cachedMetrics, err := redisCache.Get(ctx, city)
	require.NoError(t, err)
	assert.Nil(t, cachedMetrics)
}

func TestRedisCache_Delete(t *testing.T) {
	// Skip if Redis is not available
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0, 10*time.Minute)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			t.Logf("Failed to close Redis cache: %v", err)
		}
	}()

	ctx := context.Background()
	city := "Paris"
	metrics := weather.Metrics{
		Temperature: 18.0,
		Humidity:    70.0,
		Description: "Sunny",
		City:        city,
	}

	// Set data
	err = redisCache.Set(ctx, city, metrics)
	require.NoError(t, err)

	// Verify it exists
	cachedMetrics, err := redisCache.Get(ctx, city)
	require.NoError(t, err)
	require.NotNil(t, cachedMetrics)

	// Delete data
	err = redisCache.Delete(ctx, city)
	require.NoError(t, err)

	// Verify it's gone
	cachedMetrics, err = redisCache.Get(ctx, city)
	require.NoError(t, err)
	assert.Nil(t, cachedMetrics)
}

func TestRedisCache_TTL(t *testing.T) {
	// Skip if Redis is not available
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0, 100*time.Millisecond)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			t.Logf("Failed to close Redis cache: %v", err)
		}
	}()

	ctx := context.Background()
	city := "Tokyo"
	metrics := weather.Metrics{
		Temperature: 25.0,
		Humidity:    80.0,
		Description: "Rainy",
		City:        city,
	}

	// Set data with short TTL
	err = redisCache.Set(ctx, city, metrics)
	require.NoError(t, err)

	// Verify it exists immediately
	cachedMetrics, err := redisCache.Get(ctx, city)
	require.NoError(t, err)
	require.NotNil(t, cachedMetrics)

	// Wait for TTL to expire
	time.Sleep(200 * time.Millisecond)

	// Verify it's expired
	cachedMetrics, err = redisCache.Get(ctx, city)
	require.NoError(t, err)
	assert.Nil(t, cachedMetrics)
}
