package handlers

import (
	weather "Weather-Forecast-API/internal/external/openweather"
	"Weather-Forecast-API/internal/utilities"
	"context"
	"net/http"
	"time"
)

type WeatherHandler struct {
	provider weather.OpenWeatherProvider
	timeout  time.Duration
}

func NewWeatherHandler(provider weather.OpenWeatherProvider) *WeatherHandler {
	return &WeatherHandler{
		provider: provider,
		timeout:  5 * time.Second,
	}
}
func (h *WeatherHandler) GetWeather(writer http.ResponseWriter, request *http.Request) {
	city := request.URL.Query().Get("city")
	if city == "" {
		utilities.RespondJSON(writer, http.StatusBadRequest, "City parameter is required")
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), h.timeout)
	defer cancel()

	data, err := h.provider.GetWeatherByCity(ctx, city)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		return
	}

	utilities.RespondDataJSON(writer, http.StatusOK, data)

}