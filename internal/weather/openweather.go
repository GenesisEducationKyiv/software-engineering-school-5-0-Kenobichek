package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OpenWeather struct {
	APIKey string
}

func (ow OpenWeather) GetWeather(city string) (WeatherData, error) {
	geo_url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", url.QueryEscape(city), ow.APIKey)
	resp, err := http.Get(geo_url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("failed to fetch geo data")
	}
	defer resp.Body.Close()

	var geo []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return WeatherData{}, fmt.Errorf("invalid geo response")
	}

	if len(geo) == 0 {
		return WeatherData{}, fmt.Errorf("city not found")
	}

	lat, ok1 := geo[0]["lat"].(float64)
	lon, ok2 := geo[0]["lon"].(float64)
	if !ok1 || !ok2 {
		return WeatherData{}, fmt.Errorf("invalid coordinates")
	}

	weather_url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lon, ow.APIKey)
	weather_resp, err := http.Get(weather_url)
	if err != nil || weather_resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("failed to fetch weather data")
	}
	defer weather_resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(weather_resp.Body).Decode(&data); err != nil {
		return WeatherData{}, err
	}

	main, _ := data["main"].(map[string]interface{})
	wList, _ := data["weather"].([]interface{})
	wItem, _ := wList[0].(map[string]interface{})

	return WeatherData{
		Temperature: main["temp"].(float64),
		Humidity:    main["humidity"].(float64),
		Description: wItem["description"].(string),
	}, nil
}
