package weather

import (
	"Weather-Forecast-API/internal/response"
	"Weather-Forecast-API/internal/weather"
	"context"
	"net/http"
	"os"
	"time"
)

type WeatherHandler struct {
	provider weather.WeatherProvider
	timeout  time.Duration
}

func NewWeatherHandler(provider weather.WeatherProvider) *WeatherHandler {
	return &WeatherHandler{
		provider: provider,
		timeout:  5 * time.Second,
	}
}

func NewWeatherHandlerWithDefault() *WeatherHandler {
	provider := weather.NewOpenWeatherProvider(os.Getenv("OPENWEATHERMAP_API_KEY"))
	return NewWeatherHandler(provider)
}

func (h *WeatherHandler) GetWeather(writer http.ResponseWriter, request *http.Request) {
	city := request.URL.Query().Get("city")
	if city == "" {
		response.RespondJSON(writer, http.StatusBadRequest, "City parameter is required")
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), h.timeout)
	defer cancel()

	data, err := h.provider.GetWeatherByCity(ctx, city)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		return
	}

	response.RespondDataJSON(writer, http.StatusOK, data)

}
