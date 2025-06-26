package weatherapi_test

import (
	"Weather-Forecast-API/external/weatherapi"
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewWeatherAPIProvider(t *testing.T) {
	client := &http.Client{}
	apiURL := "https://api.weatherapi.com/v1/current.json"
	apiKey := "test-key"

	provider := weatherapi.NewWeatherAPIProvider(client, apiURL, apiKey)

	if provider == nil {
		t.Fatal("expected provider to be created, got nil")
	}

	// Test that the provider can be used (indirect test of field initialization)
	// We'll test the actual functionality in other tests
}

func TestGetWeather(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		apiResponse string
		apiStatus   int
		expected    weather.Metrics
		expectErr   bool
	}{
		{
			name: "valid response",
			city: "Moscow",
			apiResponse: `{
				"location": {
					"name": "Moscow"
				},
				"current": {
					"temp_c": 15.5,
					"humidity": 65,
					"condition": {
						"text": "Partly cloudy"
					}
				}
			}`,
			apiStatus: http.StatusOK,
			expected: weather.Metrics{
				Temperature: 15.5,
				Humidity:    65,
				Description: "Partly cloudy",
				City:        "Moscow",
			},
			expectErr: false,
		},
		{
			name:        "invalid JSON response",
			city:        "Moscow",
			apiResponse: `invalid-json`,
			apiStatus:   http.StatusOK,
			expected:    weather.Metrics{},
			expectErr:   true,
		},
		{
			name:        "API error response",
			city:        "Moscow",
			apiResponse: `{"error":{"code":1006,"message":"No matching location found."}}`,
			apiStatus:   http.StatusBadRequest,
			expected:    weather.Metrics{},
			expectErr:   true,
		},
		{
			name: "missing location field in API response",
			city: "Moscow",
			apiResponse: `{
				"current": {
					"temp_c": 15.5,
					"humidity": 65,
					"condition": {
						"text": "Partly cloudy"
					}
				}
			}`,
			apiStatus: http.StatusOK,
			expected: weather.Metrics{
				Temperature: 15.5,
				Humidity:    65,
				Description: "Partly cloudy",
				City:        "",
			},
			expectErr: false,
		},
		{
			name: "missing current field in API response",
			city: "Moscow",
			apiResponse: `{
				"location": {
					"name": "Moscow"
				}
			}`,
			apiStatus: http.StatusOK,
			expected: weather.Metrics{
				Temperature: 0,
				Humidity:    0,
				Description: "",
				City:        "London",
			},
			expectErr: false,
		},
		{
			name: "missing condition field in API response",
			city: "London",
			apiResponse: `{
				"location": {
					"name": "London"
				},
				"current": {
					"temp_c": 15.5,
					"humidity": 65
				}
			}`,
			apiStatus: http.StatusOK,
			expected: weather.Metrics{
				Temperature: 15.5,
				Humidity:    65,
				Description: "",
				City:        "London",
			},
			expectErr: false,
		},
		{
			name: "zero values in response",
			city: "London",
			apiResponse: `{
				"location": {
					"name": "London"
				},
				"current": {
					"temp_c": 0,
					"humidity": 0,
					"condition": {
						"text": ""
					}
				}
			}`,
			apiStatus: http.StatusOK,
			expected: weather.Metrics{
				Temperature: 0,
				Humidity:    0,
				Description: "",
				City:        "London",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request URL contains the expected parameters
				expectedURL := "/?key=test-key&q=" + tt.city
				if r.URL.String() != expectedURL {
					t.Errorf("expected URL %s, got %s", expectedURL, r.URL.String())
				}

				w.WriteHeader(tt.apiStatus)
				if _, err := w.Write([]byte(tt.apiResponse)); err != nil {
					t.Fatalf("failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client := server.Client()
			api := weatherapi.NewWeatherAPIProvider(client, server.URL, "test-key")
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			result, err := api.GetWeather(ctx, tt.city)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %+v, got: %+v", tt.expected, result)
			}
		})
	}
}

func TestGetWeather_NetworkError(t *testing.T) {
	// Create a provider with an invalid URL to simulate network error
	client := &http.Client{}
	api := weatherapi.NewWeatherAPIProvider(client, "http://invalid-url-that-does-not-exist", "test-key")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := api.GetWeather(ctx, "Moscow")
	if err == nil {
		t.Error("expected error for network failure, got nil")
	}
}

func TestGetWeather_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		response := `{
			"location": {"name": "Moscow"},
			"current": {
				"temp_c": 15.5,
				"humidity": 65,
				"condition": {"text": "Partly cloudy"}
			}
		}`
		if _, err := w.Write([]byte(response)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := server.Client()
	api := weatherapi.NewWeatherAPIProvider(client, server.URL, "test-key")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := api.GetWeather(ctx, "Moscow")
	if err == nil {
		t.Error("expected error for context cancellation, got nil")
	}
}
