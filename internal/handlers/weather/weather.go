package weather

import (
	"Weather-Forecast-API/internal/response"
	"Weather-Forecast-API/internal/weather_provider"
	"context"
	"net/http"
	"time"
)

type WeatherManager interface {
	GetWeather(writer http.ResponseWriter, request *http.Request)
}

type WeatherHandler struct {
	weatherProvider weather_provider.WeatherProvider
	requestTimeout  time.Duration
}

func NewWeatherHandler(provider weather_provider.WeatherProvider, timeout time.Duration) WeatherHandler {
	return WeatherHandler{
		weatherProvider: provider,
		requestTimeout:  timeout,
	}
}

func (h *WeatherHandler) GetWeather(writer http.ResponseWriter, request *http.Request) {
	city := request.URL.Query().Get("city")
	if city == "" {
		response.RespondJSON(writer, http.StatusBadRequest, "City parameter is required")
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), h.requestTimeout)
	defer cancel()

	data, err := h.weatherProvider.GetWeatherByCity(ctx, city)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		return
	}

	response.RespondDataJSON(writer, http.StatusOK, data)

}
