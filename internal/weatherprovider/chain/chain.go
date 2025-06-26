package chain

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"fmt"
	"log"
)

type weatherProvider interface {
	GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error)
}

type weatherChainHandler interface {
	GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error)
	SetNext(next weatherChainHandler)
}

type ChainWeatherProvider struct {
	provider weatherProvider
	next     weatherChainHandler
}

func (c *ChainWeatherProvider) SetNext(next weatherChainHandler) {
	c.next = next
}


func NewChainOpenWeatherProvider(provider weatherProvider) *ChainWeatherProvider {
	return &ChainWeatherProvider{
		provider: provider,
	}
}

func (c *ChainWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	metrics, err := c.provider.GetWeatherByCity(ctx, city)
	if err == nil {
		return metrics, nil
	}

	log.Printf("OpenWeather provider failed: %v, trying next provider", err)

	if c.next != nil {
		return c.next.GetWeatherByCity(ctx, city)
	}

	return weather.Metrics{}, fmt.Errorf("OpenWeather provider failed and no fallback available: %w", err)
}
