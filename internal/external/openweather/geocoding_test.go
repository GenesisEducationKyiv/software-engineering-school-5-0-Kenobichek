package openweather_test

import (
	"Weather-Forecast-API/internal/external/openweather"
	"Weather-Forecast-API/internal/weather/models"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGeocodingService_GetCoordinates(t *testing.T) {
	tests := []struct {
		name         string
		city         string
		mockResponse string
		mockStatus   int
		expectError  bool
		expected     models.Coordinates
	}{
		{
			name:         "valid city",
			city:         "Paris",
			mockResponse: `[{"lat": 48.8566, "lon": 2.3522}]`,
			mockStatus:   http.StatusOK,
			expectError:  false,
			expected:     models.Coordinates{Lat: 48.8566, Lon: 2.3522},
		},
		{
			name:         "city not found",
			city:         "UnknownCity",
			mockResponse: `[]`,
			mockStatus:   http.StatusOK,
			expectError:  true,
			expected:     models.Coordinates{},
		},
		{
			name:         "API returns error status",
			city:         "InvalidCity",
			mockResponse: `{"message":"error"}`,
			mockStatus:   http.StatusBadRequest,
			expectError:  true,
			expected:     models.Coordinates{},
		},
		{
			name:         "invalid JSON response",
			city:         "CorruptDataCity",
			mockResponse: `{"message":"invalid json"}`,
			mockStatus:   http.StatusOK,
			expectError:  true,
			expected:     models.Coordinates{},
		},
		{
			name:         "empty city name",
			city:         "",
			mockResponse: ``,
			mockStatus:   http.StatusOK,
			expectError:  true,
			expected:     models.Coordinates{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedCity := r.URL.Query().Get("q")
				if receivedCity != url.QueryEscape(tt.city) {
					t.Errorf("expected city %q, got %q", tt.city, receivedCity)
				}
				w.WriteHeader(tt.mockStatus)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			service := &openweather.GeocodingService{
				ApiKey:  "testAPIKey",
				BaseURL: mockServer.URL,
			}

			result, err := service.GetCoordinates(context.Background(), tt.city)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, got %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
