package openweather_test

import (
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWeather(t *testing.T) {
	tests := []struct {
		name        string
		coords      weather.Coordinates
		apiResponse string
		apiStatus   int
		expected    weather.Metrics
		expectErr   bool
	}{
		{
			name: "valid response",
			coords: weather.Coordinates{
				Lat: 37.7749,
				Lon: -122.4194,
			},
			apiResponse: `{"main":{"temp":20.5,"humidity":50},"weather":[{"description":"clear sky"}]}`,
			apiStatus:   http.StatusOK,
			expected: weather.Metrics{
				Temperature: 20.5,
				Humidity:    50,
				Description: "clear sky",
			},
			expectErr: false,
		},
		{
			name: "invalid JSON response",
			coords: weather.Coordinates{
				Lat: 37.7749,
				Lon: -122.4194,
			},
			apiResponse: `invalid-json`,
			apiStatus:   http.StatusOK,
			expected:    weather.Metrics{},
			expectErr:   true,
		},
		{
			name: "API error response",
			coords: weather.Coordinates{
				Lat: 37.7749,
				Lon: -122.4194,
			},
			apiResponse: `{"error":"something went wrong"}`,
			apiStatus:   http.StatusInternalServerError,
			expected:    weather.Metrics{},
			expectErr:   true,
		},
		{
			name: "missing fields in API response",
			coords: weather.Coordinates{
				Lat: 37.7749,
				Lon: -122.4194,
			},
			apiResponse: `{"main":{"temp":20.5},"weather":[]}`,
			apiStatus:   http.StatusOK,
			expected:    weather.Metrics{},
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.apiStatus)
				if _, err := w.Write([]byte(tt.apiResponse)); err != nil {
					t.Fatalf("failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client := server.Client()
			api := openweather.NewOpenWeatherAPI(client, server.URL, "dummy-key")
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			result, err := api.GetWeather(ctx, tt.coords)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
