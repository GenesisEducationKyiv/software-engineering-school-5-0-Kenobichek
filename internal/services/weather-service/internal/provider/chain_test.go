package provider_test

import (
	"context"
	"errors"
	"testing"

	"internal/services/weather-service/internal/domain"
	provider "internal/services/weather-service/internal/provider"
)

type mockWeatherProvider struct {
	metrics domain.Metrics
	err     error
	called  *bool
}

func (m *mockWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (domain.Metrics, error) {
	if m.called != nil {
		*m.called = true
	}
	return m.metrics, m.err
}

func TestChainWeatherProvider_FirstFails_SecondUsed(t *testing.T) {
	firstCalled, secondCalled := false, false
	first := &mockWeatherProvider{err: errors.New("fail"), called: &firstCalled}
	second := &mockWeatherProvider{metrics: domain.Metrics{City: "Kyiv"}, called: &secondCalled}
	chain1 := provider.NewChainWeatherProvider(first)
	chain2 := provider.NewChainWeatherProvider(second)
	chain1.SetNext(chain2)

	result, err := chain1.GetWeatherByCity(context.Background(), "Kyiv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.City != "Kyiv" || !firstCalled || !secondCalled {
		t.Errorf("chain logic failed: %+v, firstCalled=%v, secondCalled=%v", result, firstCalled, secondCalled)
	}
}

func TestChainWeatherProvider_FirstSuccess_SecondNotUsed(t *testing.T) {
	firstCalled, secondCalled := false, false
	first := &mockWeatherProvider{metrics: domain.Metrics{City: "Kyiv"}, called: &firstCalled}
	second := &mockWeatherProvider{metrics: domain.Metrics{City: "Other"}, called: &secondCalled}
	chain1 := provider.NewChainWeatherProvider(first)
	chain2 := provider.NewChainWeatherProvider(second)
	chain1.SetNext(chain2)

	result, err := chain1.GetWeatherByCity(context.Background(), "Kyiv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.City != "Kyiv" || !firstCalled || secondCalled {
		t.Errorf("chain logic failed: %+v, firstCalled=%v, secondCalled=%v", result, firstCalled, secondCalled)
	}
}
