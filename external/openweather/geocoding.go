package openweather

import (
	"Weather-Forecast-API/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GeocodingProvider interface {
	GetCoordinates(ctx context.Context, city string) (Coordinates, error)
}

type OpenWeatherGeocodingService struct {
	cfg *config.Config
}

func NewOpenWeatherGeocodingService(cfg *config.Config) OpenWeatherGeocodingService {
	return OpenWeatherGeocodingService{
		cfg: cfg,
	}
}

func (g *OpenWeatherGeocodingService) GetCoordinates(ctx context.Context, city string) (Coordinates, error) {
	geoURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s",
		g.cfg.OpenWeather.GeocodingAPIURL, url.QueryEscape(city), g.cfg.OpenWeather.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, nil)
	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var geo []Coordinates
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return Coordinates{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geo) == 0 {
		return Coordinates{}, fmt.Errorf("city not found: %s", city)
	}

	return geo[0], nil
}
