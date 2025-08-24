package provider

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"internal/services/weather-service/internal/domain"
)

type weatherCacheManager interface {
	Get(ctx context.Context, city string) (*domain.Metrics, error)
	Set(ctx context.Context, city string, metrics domain.Metrics) error
	Delete(ctx context.Context, city string) error
	Close() error
}

type eventPublishingManager interface {
	PublishWeatherUpdated(city string, metrics domain.Metrics) error
}

type CachedWeatherProvider struct {
	provider       weatherProviderManager
	cache          weatherCacheManager
	eventPublisher eventPublishingManager
	wg             sync.WaitGroup
}

func NewCachedWeatherProvider(provider weatherProviderManager, cache weatherCacheManager) *CachedWeatherProvider {
	return &CachedWeatherProvider{
		provider: provider,
		cache:    cache,
	}
}

func NewCachedWeatherProviderWithEvents(
	provider weatherProviderManager,
	cache weatherCacheManager,
	eventPublisher eventPublishingManager,
) *CachedWeatherProvider {
	return &CachedWeatherProvider{
		provider:       provider,
		cache:          cache,
		eventPublisher: eventPublisher,
	}
}

func (c *CachedWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error) {
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
		return domain.Metrics{}, fmt.Errorf("failed to get weather from provider: %w", err)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.cache.Set(ctx, city, metrics); err != nil {
			log.Printf("Failed to cache weather data for city %s: %v", city, err)
		} else if c.eventPublisher != nil {
			if err := c.eventPublisher.PublishWeatherUpdated(city, metrics); err != nil {
				log.Printf("Failed to publish weather updated event for city %s: %v", city, err)
			}
		}
	}()

	return metrics, nil
}

func (c *CachedWeatherProvider) Close() error {
	c.wg.Wait()
	return c.cache.Close()
}
