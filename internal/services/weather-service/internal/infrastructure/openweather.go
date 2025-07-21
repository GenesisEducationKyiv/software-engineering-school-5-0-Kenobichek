package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"internal/services/weather-service/internal/domain"
)

type GeocodingService struct {
	httpClient *http.Client
	apiurl     string
	apikey     string
}

func NewGeocodingService(httpClient *http.Client, apiurl, apikey string) *GeocodingService {
	return &GeocodingService{
		httpClient: httpClient,
		apiurl:     apiurl,
		apikey:     apikey,
	}
}

func (g *GeocodingService) GetCoordinates(ctx context.Context, city string) (domain.Coordinates, error) {
	geoURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s", g.apiurl, url.QueryEscape(city), g.apikey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, http.NoBody)
	if err != nil {
		return domain.Coordinates{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return domain.Coordinates{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return domain.Coordinates{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var geo []domain.Coordinates
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return domain.Coordinates{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geo) == 0 {
		return domain.Coordinates{}, fmt.Errorf("city not found: %s", city)
	}

	return geo[0], nil
}

type OpenWeatherAPI struct {
	httpClient *http.Client
	apiurl     string
	apikey     string
}

func NewOpenWeatherAPI(httpClient *http.Client, apiurl, apikey string) *OpenWeatherAPI {
	return &OpenWeatherAPI{
		httpClient: httpClient,
		apiurl:     apiurl,
		apikey:     apikey,
	}
}

func (w *OpenWeatherAPI) GetWeather(ctx context.Context, coords domain.Coordinates) (domain.Metrics, error) {
	weatherURL := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", w.apiurl, coords.Lat, coords.Lon, w.apikey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, http.NoBody)
	if err != nil {
		return domain.Metrics{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return domain.Metrics{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return domain.Metrics{}, fmt.Errorf("API returned status code: %d", resp.StatusCode)
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
		return domain.Metrics{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.Weather) == 0 {
		return domain.Metrics{}, fmt.Errorf("no weather data available")
	}

	return domain.Metrics{
		Temperature: data.Main.Temperature,
		Humidity:    data.Main.Humidity,
		Description: data.Weather[0].Description,
	}, nil
}
