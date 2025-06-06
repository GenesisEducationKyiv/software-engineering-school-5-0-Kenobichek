package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type OpenWeather struct {
	APIKey string
}

func (ow OpenWeather) GetWeather(ctx context.Context, city string) (DataWeather, error) {
	geoURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
		url.QueryEscape(city), ow.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, nil)
	if err != nil {
		return DataWeather{}, errors.New("new request error")
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return DataWeather{}, errors.New("failed to fetch geo data")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	var geo []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return DataWeather{}, errors.New("invalid geo response")
	}

	if len(geo) == 0 {
		return DataWeather{}, errors.New("city not found")
	}

	lat, ok1 := geo[0]["lat"].(float64)
	lon, ok2 := geo[0]["lon"].(float64)

	if !ok1 || !ok2 {
		return DataWeather{}, errors.New("invalid coordinates")
	}

	weatherURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric",
		lat, lon, ow.APIKey)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, nil)
	if err != nil {
		return DataWeather{}, errors.New("new request error")
	}

	weatherResp, err := http.DefaultClient.Do(req)

	if err != nil || weatherResp.StatusCode != http.StatusOK {
		return DataWeather{}, errors.New("failed to fetch weather data")
	}

	defer func() {
		if err := weatherResp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	var data map[string]interface{}
	if err := json.NewDecoder(weatherResp.Body).Decode(&data); err != nil {
		return DataWeather{}, fmt.Errorf("failed to decode weather API response: %w", err)
	}

	main, _ := data["main"].(map[string]interface{})
	wList, _ := data["weather"].([]interface{})
	wItem, _ := wList[0].(map[string]interface{})

	return DataWeather{
		Temperature: main["temp"].(float64),
		Humidity:    main["humidity"].(float64),
		Description: wItem["description"].(string),
	}, nil
}
