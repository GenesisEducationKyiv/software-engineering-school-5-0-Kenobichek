package cached

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"fmt"
	"log"
)

type weatherChainHandler interface {
	GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error)
}

type weatherCacheManager interface {
	Get(ctx context.Context, city string) (*weather.Metrics, error)
	Set(ctx context.Context, city string, metrics weather.Metrics) error
	Delete(ctx context.Context, city string) error
	Close() error
}


type CachedWeatherProvider struct {
	provider weatherChainHandler
	cache    weatherCacheManager
}

func NewCachedWeatherProvider(provider weatherChainHandler, cache weatherCacheManager) *CachedWeatherProvider {
	return &CachedWeatherProvider{
		provider: provider,
		cache:    cache,
	}
}

func (c *CachedWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	cachedMetrics, err := c.cache.Get(ctx, city)
	if err != nil {
		log.Printf("Cache get error for city %s: %v", city, err)
	} else if cachedMetrics != nil {
		log.Printf("Cache hit for city: %s", city)
		return *cachedMetrics, nil
	}

	log.Printf("Cache miss for city: %s, fetching from provider", city)
	metrics, err := c.provider.GetWeatherByCity(ctx, city)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to get weather from provider: %w", err)
	}

	go func() {
		if err := c.cache.Set(context.Background(), city, metrics); err != nil {
			log.Printf("Failed to cache weather data for city %s: %v", city, err)
		}
	}()

	return metrics, nil
}
