package openweatherprovider_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    weatherprovider "github.com/yourorg/yourapp/internal/weatherprovider"
    owp "github.com/yourorg/yourapp/internal/weatherprovider/openweatherprovider"
)

type mockGeocodingManager struct {
    getCoordinatesFunc func(context.Context, string) (weatherprovider.Coordinates, error)
}

func (m *mockGeocodingManager) GetCoordinates(ctx context.Context, city string) (weatherprovider.Coordinates, error) {
    if m.getCoordinatesFunc != nil {
        return m.getCoordinatesFunc(ctx, city)
    }
    return weatherprovider.Coordinates{}, nil
}

type mockWeatherManager struct {
    getWeatherFunc func(context.Context, weatherprovider.Coordinates) (weatherprovider.Metrics, error)
}

func (m *mockWeatherManager) GetWeather(ctx context.Context, coords weatherprovider.Coordinates) (weatherprovider.Metrics, error) {
    if m.getWeatherFunc != nil {
        return m.getWeatherFunc(ctx, coords)
    }
    return weatherprovider.Metrics{}, nil
}

func TestOpenWeatherProvider_GetWeatherByCity(t *testing.T) {
    fixedCoords := weatherprovider.Coordinates{Lat: 40.7128, Lon: -74.0060}
    fixedMetrics := weatherprovider.Metrics{Temperature: 22.5, Humidity: 60, Pressure: 1013}

    tests := []struct {
        name               string
        ctx                context.Context
        city               string
        geoMock            *mockGeocodingManager
        weatherMock        *mockWeatherManager
        expectedMetrics    weatherprovider.Metrics
        expectedErr        error
    }{
        {
            name: "successful retrieval",
            ctx:  context.Background(),
            city: "New York",
            geoMock: &mockGeocodingManager{
                getCoordinatesFunc: func(_ context.Context, city string) (weatherprovider.Coordinates, error) {
                    assert.Equal(t, "New York", city)
                    return fixedCoords, nil
                },
            },
            weatherMock: &mockWeatherManager{
                getWeatherFunc: func(_ context.Context, coords weatherprovider.Coordinates) (weatherprovider.Metrics, error) {
                    assert.Equal(t, fixedCoords, coords)
                    return fixedMetrics, nil
                },
            },
            expectedMetrics: fixedMetrics,
            expectedErr:     nil,
        },
        {
            name:    "geocoding error",
            ctx:     context.Background(),
            city:    "Nowhere",
            geoMock: &mockGeocodingManager{getCoordinatesFunc: func(_ context.Context, _ string) (weatherprovider.Coordinates, error) {
                return weatherprovider.Coordinates{}, errors.New("city not found")
            }},
            weatherMock:     &mockWeatherManager{},
            expectedMetrics: weatherprovider.Metrics{},
            expectedErr:     errors.New("city not found"),
        },
        {
            name:    "weather API error",
            ctx:     context.Background(),
            city:    "Rome",
            geoMock: &mockGeocodingManager{getCoordinatesFunc: func(_ context.Context, _ string) (weatherprovider.Coordinates, error) {
                return fixedCoords, nil
            }},
            weatherMock: &mockWeatherManager{getWeatherFunc: func(_ context.Context, _ weatherprovider.Coordinates) (weatherprovider.Metrics, error) {
                return weatherprovider.Metrics{}, errors.New("upstream service failure")
            }},
            expectedMetrics: weatherprovider.Metrics{},
            expectedErr:     errors.New("upstream service failure"),
        },
        {
            name:            "empty city validation",
            ctx:             context.Background(),
            city:            "",
            geoMock:         &mockGeocodingManager{},
            weatherMock:     &mockWeatherManager{},
            expectedMetrics: weatherprovider.Metrics{},
            expectedErr:     errors.New("city cannot be empty"),
        },
        {
            name:            "nil context validation",
            ctx:             nil,
            city:            "Paris",
            geoMock:         &mockGeocodingManager{},
            weatherMock:     &mockWeatherManager{},
            expectedMetrics: weatherprovider.Metrics{},
            expectedErr:     errors.New("context is required"),
        },
        {
            name: "context timeout",
            ctx: func() context.Context {
                c, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
                defer cancel()
                return c
            }(),
            city: "Berlin",
            geoMock: &mockGeocodingManager{getCoordinatesFunc: func(ctx context.Context, _ string) (weatherprovider.Coordinates, error) {
                // exceed deadline
                time.Sleep(20 * time.Millisecond)
                return weatherprovider.Coordinates{}, nil
            }},
            weatherMock:     &mockWeatherManager{},
            expectedMetrics: weatherprovider.Metrics{},
            expectedErr:     context.DeadlineExceeded,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            provider := owp.New(tc.geoMock, tc.weatherMock)
            metrics, err := provider.GetWeatherByCity(tc.ctx, tc.city)

            if tc.expectedErr != nil {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tc.expectedErr.Error())
                assert.Empty(t, metrics)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedMetrics, metrics)
            }
        })
    }
}