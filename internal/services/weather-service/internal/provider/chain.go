package provider

import (
	"context"
	"fmt"
	"log"

	"internal/services/weather-service/internal/domain"
)

type weatherProviderManager interface {
	GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error)
}

type WeatherChainHandler interface {
	GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error)
	SetNext(next WeatherChainHandler)
}

type ChainWeatherProvider struct {
	provider weatherProviderManager
	next     WeatherChainHandler
}

func NewChainWeatherProvider(provider weatherProviderManager) *ChainWeatherProvider {
	return &ChainWeatherProvider{
		provider: provider,
	}
}

func (c *ChainWeatherProvider) SetNext(next WeatherChainHandler) {
	c.next = next
}

func (c *ChainWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error) {
	metrics, err := c.provider.GetWeatherByCity(ctx, city)
	if err == nil {
		return metrics, nil
	}

	log.Printf("Weather provider failed: %v, trying next provider", err)

	if c.next != nil {
		return c.next.GetWeatherByCity(ctx, city)
	}

	return domain.Metrics{}, fmt.Errorf("no fallback provider: %w", err)
}
