package openweather_test

import (
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeocodingService_GetCoordinates(t *testing.T) {
	tests := []struct {
		name         string
		city         string
		mockResponse string
		mockStatus   int
		expected     weather.Coordinates
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "Valid City",
			city:         "Paris",
			mockResponse: `[{"lat": 48.8566, "lon": 2.3522}]`,
			mockStatus:   http.StatusOK,
			expected:     weather.Coordinates{Lat: 48.8566, Lon: 2.3522},
			expectError:  false,
		},
		{
			name:         "City Not Found",
			city:         "InvalidCity",
			mockResponse: `[]`,
			mockStatus:   http.StatusOK,
			expected:     weather.Coordinates{},
			expectError:  true,
			errorMsg:     "city not found: InvalidCity",
		},
		{
			name:         "Invalid JSON Response",
			city:         "Paris",
			mockResponse: `{"invalid"`,
			mockStatus:   http.StatusOK,
			expected:     weather.Coordinates{},
			expectError:  true,
			errorMsg:     "failed to decode response",
		},
		{
			name:         "Non-200 Status Code",
			city:         "Paris",
			mockResponse: `{"message": "error"}`,
			mockStatus:   http.StatusInternalServerError,
			expected:     weather.Coordinates{},
			expectError:  true,
			errorMsg:     "API returned status code: 500",
		},
		{
			name:         "Empty Response Body",
			city:         "Paris",
			mockResponse: "",
			mockStatus:   http.StatusOK,
			expected:     weather.Coordinates{},
			expectError:  true,
			errorMsg:     "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				if tt.mockResponse != "" {
					_, _ = w.Write([]byte(tt.mockResponse))
				}
			}))
			defer server.Close()

			// Create a new GeocodingService instance
			service := openweather.NewGeocodingService(server.Client(), server.URL, "test-api-key")

			// Call GetCoordinates
			ctx := context.Background()
			result, err := service.GetCoordinates(ctx, tt.city)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
