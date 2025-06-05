package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type OpenWeather struct {
	APIKey string
}

func (ow OpenWeather) GetWeather(city string) (DataWeather, error) {
	geoURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
		url.QueryEscape(city), ow.APIKey)
	resp, err := http.Get(geoURL)

	if err != nil || resp.StatusCode != http.StatusOK {
		return DataWeather{}, errors.New("failed to fetch geo data")
	}
	defer resp.Body.Close()

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

	weatherResp, err := http.Get(weatherURL)
	if err != nil || weatherResp.StatusCode != http.StatusOK {
		return DataWeather{}, errors.New("failed to fetch weather data")
	}
	defer weatherResp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(weatherResp.Body).Decode(&data); err != nil {
		return DataWeather{}, err
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
