package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api-gateway/proto"
)

const (
	defaultRequestTimeout = 5 * time.Second
)

type weatherClientManager interface {
	GetWeather(ctx context.Context, req *proto.WeatherRequest) (*proto.WeatherResponse, error)
}

type WeatherHandler struct {
	weatherClient weatherClientManager
}

func NewWeatherHandler(weatherClient weatherClientManager) *WeatherHandler {
	return &WeatherHandler{weatherClient: weatherClient}
}

func (h *WeatherHandler) WeatherProxyHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	log.Printf("[WeatherProxyHandler] incoming request: %s %s, city=%s", r.Method, r.URL.Path, city)
	if err := validateWeatherParams(city); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("[WeatherProxyHandler] missing city parameter")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
	defer cancel()
	resp, err := h.weatherClient.GetWeather(ctx, &proto.WeatherRequest{City: city})
	if err != nil {
		http.Error(w, "failed to get weather: "+err.Error(), http.StatusBadGateway)
		log.Printf("[WeatherProxyHandler] gRPC error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[WeatherProxyHandler] failed to encode response: %v", err)
	}
	log.Printf("[WeatherProxyHandler] success for city=%s", city)
}
