package openweather

import (
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
	httpClient *http.Client
	apiurl     string
	apikey     string
}

func NewOpenWeatherGeocodingService(
	httpClient *http.Client,
	apiurl string,
	apikey string) *OpenWeatherGeocodingService {
	return &OpenWeatherGeocodingService{
		httpClient: httpClient,
		apiurl:     apiurl,
		apikey:     apikey,
	}
}

func (g *OpenWeatherGeocodingService) GetCoordinates(ctx context.Context, city string) (Coordinates, error) {
	geoURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s",
		g.apiurl, url.QueryEscape(city), g.apikey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, http.NoBody)
	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
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
