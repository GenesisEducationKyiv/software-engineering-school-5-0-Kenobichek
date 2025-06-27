package chain

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"errors"
	"fmt"
	"log"
)

var ErrNoFallback = errors.New("weather provider failed and no fallback available")

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

func NewChainWeatherProvider(provider weatherProvider) *ChainWeatherProvider {
	return &ChainWeatherProvider{
		provider: provider,
	}
}

func (c *ChainWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error) {
	metrics, err := c.provider.GetWeatherByCity(ctx, city)
	if err == nil {
		return metrics, nil
	}

	log.Printf("Weather provider failed: %v, trying next provider", err)

	if c.next != nil {
		return c.next.GetWeatherByCity(ctx, city)
	}

	return weather.Metrics{}, fmt.Errorf("%s: %w", ErrNoFallback, err)
}
