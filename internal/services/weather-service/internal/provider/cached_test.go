package provider_test

import (
	"context"
	"testing"

	"internal/services/weather-service/internal/domain"
	provider "internal/services/weather-service/internal/provider"
)

type mockCache struct {
	metrics *domain.Metrics
	hit     bool
}

func (m *mockCache) Get(ctx context.Context, city string) (*domain.Metrics, error) {
	if m.hit {
		return m.metrics, nil
	}
	return nil, nil
}
func (m *mockCache) Set(ctx context.Context, city string, metrics domain.Metrics) error {
	m.metrics = &metrics
	m.hit = true
	return nil
}
func (m *mockCache) Delete(ctx context.Context, city string) error { return nil }
func (m *mockCache) Close() error                                  { return nil }

type mockWeatherProviderCached struct {
	metrics domain.Metrics
}

func (m *mockWeatherProviderCached) GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error) {
	return m.metrics, nil
}

func TestCachedWeatherProvider_CacheHit(t *testing.T) {
	cache := &mockCache{metrics: &domain.Metrics{City: "Kyiv"}, hit: true}
	prov := &mockWeatherProviderCached{metrics: domain.Metrics{City: "Kyiv"}}
	cached := provider.NewCachedWeatherProvider(prov, cache)

	result, err := cached.GetWeatherByCity(context.Background(), "Kyiv")
	if err != nil || result.City != "Kyiv" {
		t.Errorf("expected cache hit, got: %+v, err: %v", result, err)
	}
}

func TestCachedWeatherProvider_CacheMiss(t *testing.T) {
	cache := &mockCache{hit: false}
	prov := &mockWeatherProviderCached{metrics: domain.Metrics{City: "Kyiv"}}
	cached := provider.NewCachedWeatherProvider(prov, cache)

	result, err := cached.GetWeatherByCity(context.Background(), "Kyiv")
	if err != nil || result.City != "Kyiv" {
		t.Errorf("expected cache miss logic, got: %+v, err: %v", result, err)
	}
	err = cached.Close()
	if err != nil {
		t.Errorf("failed to close cache: %v", err)
	}
	if !cache.hit {
		t.Errorf("expected cache to be set after miss, but cache.hit is false")
	}
}
