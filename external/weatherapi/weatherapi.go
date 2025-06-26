package weatherapi

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherAPIProvider struct {
	httpClient *http.Client
	apiurl     string
	apikey     string
}

func NewWeatherAPIProvider(httpClient *http.Client, apiurl, apikey string) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		httpClient: httpClient,
		apiurl:     apiurl,
		apikey:     apikey,
	}
}

func (w *WeatherAPIProvider) GetWeather(ctx context.Context, city string) (weather.Metrics, error) {
	weatherURL := fmt.Sprintf("%s?key=%s&q=%s", w.apiurl, w.apikey, city)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, http.NoBody)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.Metrics{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var data struct {
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
		Current struct {
			TempC     float64 `json:"temp_c"`
			Humidity  float64 `json:"humidity"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return weather.Metrics{
		Temperature: data.Current.TempC,
		Humidity:    data.Current.Humidity,
		Description: data.Current.Condition.Text,
		City:        data.Location.Name,
	}, nil
}
