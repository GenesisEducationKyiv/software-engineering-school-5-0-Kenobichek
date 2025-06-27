package weather

import (
	"Weather-Forecast-API/internal/response"
	"context"
	"net/http"
	"time"
)

type weatherProviderManager interface {
	GetWeatherByCity(ctx context.Context, city string) (Metrics, error)
}

type Handler struct {
	weatherProvider weatherProviderManager
	requestTimeout  time.Duration
}

func NewHandler(
	provider weatherProviderManager,
	timeout time.Duration,
) *Handler {
	return &Handler{
		weatherProvider: provider,
		requestTimeout:  timeout,
	}
}

func (h *Handler) GetWeather(writer http.ResponseWriter, request *http.Request) {
	city := request.URL.Query().Get("city")
	if city == "" {
		response.RespondJSON(writer, http.StatusBadRequest, "City parameter is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	data, err := h.weatherProvider.GetWeatherByCity(ctx, city)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		return
	}

	response.RespondDataJSON(writer, http.StatusOK, data)

}
