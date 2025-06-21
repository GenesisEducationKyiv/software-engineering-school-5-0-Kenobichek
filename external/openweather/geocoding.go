package openweather

import (
	"Weather-Forecast-API/internal/handlers/weather"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GeocodingService struct {
	httpClient *http.Client
	apiurl     string
	apikey     string
}

func NewGeocodingService(
	httpClient *http.Client,
	apiurl string,
	apikey string,
) *GeocodingService {
	return &GeocodingService{
		httpClient: httpClient,
		apiurl:     apiurl,
		apikey:     apikey,
	}
}

func (g *GeocodingService) GetCoordinates(ctx context.Context, city string) (weather.Coordinates, error) {
	geoURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s",
		g.apiurl, url.QueryEscape(city), g.apikey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, http.NoBody)
	if err != nil {
		return weather.Coordinates{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return weather.Coordinates{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return weather.Coordinates{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var geo []weather.Coordinates
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return weather.Coordinates{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geo) == 0 {
		return weather.Coordinates{}, fmt.Errorf("city not found: %s", city)
	}

	return geo[0], nil
}
