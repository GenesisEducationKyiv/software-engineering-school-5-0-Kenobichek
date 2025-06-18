package openweather

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/weather_provider/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenWeatherAPI struct {
	cfg *config.Config
}

func NewOpenWeatherAPI(cfg *config.Config) OpenWeatherAPI {
	return OpenWeatherAPI{
		cfg: cfg,
	}
}

func (w *OpenWeatherAPI) GetWeather(ctx context.Context, coords models.Coordinates) (models.WeatherData, error) {
	weatherURL := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric",
		w.cfg.OpenWeather.WeatherAPIURL, coords.Lat, coords.Lon, w.cfg.OpenWeather.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, nil)
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return models.WeatherData{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
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
		return models.WeatherData{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.Weather) == 0 {
		return models.WeatherData{}, fmt.Errorf("no weather data available")
	}

	weatherData := models.WeatherData{
		Temperature: data.Main.Temperature,
		Humidity:    data.Main.Humidity,
		Description: data.Weather[0].Description,
	}

	return weatherData, nil
}
