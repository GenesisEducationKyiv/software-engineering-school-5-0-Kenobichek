package weather

import (
    "context"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
)

type mockWeatherProvider struct {
    getWeatherFunc func(context.Context, string) (Metrics, error)
}

func (m *mockWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (Metrics, error) {
    if m.getWeatherFunc != nil {
        return m.getWeatherFunc(ctx, city)
    }
    return Metrics{}, nil
}

func TestHandler_GetWeather(t *testing.T) {
    tests := []struct {
        name           string
        city           string
        weatherData    Metrics
        weatherError   error
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "successful weather retrieval",
            city: "London",
            weatherData: Metrics{
                Temperature: 20.5,
                Humidity:    65.0,
                Description: "Partly cloudy",
            },
            expectedStatus: http.StatusOK,
            expectedBody:   `{"data":{"Temperature":20.5,"Humidity":65,"Description":"Partly cloudy"}}`,
        },
        {
            name:           "missing city parameter",
            city:           "",
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"message":"City parameter is required"}`,
        },
        {
            name:           "weather service error",
            city:           "InvalidCity",
            weatherError:   errors.New("city not found"),
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"message":"Failed to get weather: city not found"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider := &mockWeatherProvider{
                getWeatherFunc: func(ctx context.Context, city string) (Metrics, error) {
                    if tt.weatherError != nil {
                        return Metrics{}, tt.weatherError
                    }
                    return tt.weatherData, nil
                },
            }

            handler := NewHandler(provider, 5*time.Second)

            url := "/weather"
            if tt.city != "" {
                url += "?city=" + tt.city
            }
            req := httptest.NewRequest(http.MethodGet, url, nil)
            w := httptest.NewRecorder()

            handler.GetWeather(w, req)

            if w.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
            }

            body := strings.TrimSpace(w.Body.String())
            if body != tt.expectedBody {
                t.Errorf("expected body %s, got %s", tt.expectedBody, body)
            }
        })
    }
}

func TestHandler_GetWeatherTimeout(t *testing.T) {
    provider := &mockWeatherProvider{
        getWeatherFunc: func(ctx context.Context, city string) (Metrics, error) {
            select {
            case <-time.After(100 * time.Millisecond):
                return Metrics{}, nil
            case <-ctx.Done():
                return Metrics{}, ctx.Err()
            }
        },
    }

    handler := NewHandler(provider, 50*time.Millisecond)

    req := httptest.NewRequest(http.MethodGet, "/weather?city=London", nil)
    w := httptest.NewRecorder()

    handler.GetWeather(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
    }

    if !strings.Contains(w.Body.String(), "context deadline exceeded") {
        t.Errorf("expected timeout error, got %s", w.Body.String())
    }
}