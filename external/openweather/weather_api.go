package openweather

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WeatherAPI struct {
	httpClient *http.Client
	apiurl     string
	apikey     string
}

func NewOpenWeatherAPI(
	httpClient *http.Client,
	apiurl string,
	apikey string,
) *WeatherAPI {
	return &WeatherAPI{
		httpClient: httpClient,
		apiurl:     apiurl,
		apikey:     apikey,
	}
}

func (w *WeatherAPI) GetWeather(ctx context.Context, coords weather.Coordinates) (weather.Metrics, error) {
	weatherURL := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric",
		w.apiurl, coords.Lat, coords.Lon, w.apikey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, http.NoBody)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return weather.Metrics{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var data struct {
		Main struct {
			Temperature float64 `json:"temp"`
			Humidity    float64 `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weather.Metrics{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.Weather) == 0 {
		return weather.Metrics{}, fmt.Errorf("no weather data available")
	}

	weatherData := weather.Metrics{
		Temperature: data.Main.Temperature,
		Humidity:    data.Main.Humidity,
		Description: data.Weather[0].Description,
	}

	return weatherData, nil
}
