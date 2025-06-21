package openweatherprovider_test

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/weatherprovider/openweatherprovider"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockGeocodingManager struct {
	GetCoordinatesFn func(ctx context.Context, city string) (weather.Coordinates, error)
}

func (m *mockGeocodingManager) GetCoordinates(ctx context.Context, city string) (weather.Coordinates, error) {
	return m.GetCoordinatesFn(ctx, city)
}

type mockWeatherManager struct {
	GetWeatherFn func(ctx context.Context, coords weather.Coordinates) (weather.Metrics, error)
}

func (m *mockWeatherManager) GetWeather(ctx context.Context, coords weather.Coordinates) (weather.Metrics, error) {
	return m.GetWeatherFn(ctx, coords)
}

func TestGetWeatherByCity(t *testing.T) {
	tests := []struct {
		name             string
		city             string
		mockCoords       weather.Coordinates
		mockCoordsErr    error
		mockWeather      weather.Metrics
		mockWeatherErr   error
		expectedMetrics  weather.Metrics
		expectedErr      string
		cancelContext    bool
		geocodingManager mockGeocodingManager
		weatherManager   mockWeatherManager
	}{
		{
			name:            "valid city",
			city:            "New York",
			mockCoords:      weather.Coordinates{Lat: 40.7128, Lon: -74.0060},
			mockWeather:     weather.Metrics{Temperature: 22.5, Humidity: 60.0, Description: "Sunny"},
			expectedMetrics: weather.Metrics{Temperature: 22.5, Humidity: 60.0, Description: "Sunny"},
		},
		{
			name:          "empty city",
			city:          "  ",
			expectedErr:   "city must not be empty",
			mockCoordsErr: errors.New("city must not be empty"),
		},
		{
			name:          "geocoding failed",
			city:          "Nonexistent",
			mockCoordsErr: errors.New("failed to get coordinates"),
			expectedErr:   "failed to get coordinates: failed to get coordinates",
		},
		{
			name:           "weather API failed",
			city:           "Paris",
			mockCoords:     weather.Coordinates{Lat: 48.8566, Lon: 2.3522},
			mockWeatherErr: errors.New("weather API unavailable"),
			expectedErr:    "failed to get weather: weather API unavailable",
		},
		{
			name:          "context cancelled",
			city:          "London",
			cancelContext: true,
			expectedErr:   "context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGeocoding := &mockGeocodingManager{
				GetCoordinatesFn: func(ctx context.Context, city string) (weather.Coordinates, error) {
					if city == tt.city {
						return tt.mockCoords, tt.mockCoordsErr
					}
					return weather.Coordinates{}, errors.New("unexpected city")
				},
			}

			mockWeather := &mockWeatherManager{
				GetWeatherFn: func(ctx context.Context, coords weather.Coordinates) (weather.Metrics, error) {
					if coords == tt.mockCoords {
						return tt.mockWeather, tt.mockWeatherErr
					}
					return weather.Metrics{}, errors.New("unexpected coordinates")
				},
			}

			ctx := context.Background()
			if tt.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			provider := openweatherprovider.NewOpenWeatherProvider(mockGeocoding, mockWeather)
			metrics, err := provider.GetWeatherByCity(ctx, tt.city)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMetrics, metrics)
			}
		})
	}
}
