package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api-gateway/proto"

	"google.golang.org/grpc"
)

type WeatherHandler struct {
	grpcClient proto.WeatherServiceClient
}

func NewWeatherHandler(grpcAddr string) (*WeatherHandler, error) {
	log.Printf("[WeatherHandler] initializing gRPC client to %s", grpcAddr)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("[WeatherHandler] failed to dial gRPC at %s: %v", grpcAddr, err)
		return nil, err
	}
	log.Printf("[WeatherHandler] gRPC connection established to %s", grpcAddr)
	client := proto.NewWeatherServiceClient(conn)
	return &WeatherHandler{grpcClient: client}, nil
}

func (h *WeatherHandler) WeatherProxyHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	log.Printf("[WeatherProxyHandler] incoming request: %s %s, city=%s", r.Method, r.URL.Path, city)
	if city == "" {
		http.Error(w, "city parameter is required", http.StatusBadRequest)
		log.Printf("[WeatherProxyHandler] missing city parameter")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := h.grpcClient.GetWeather(ctx, &proto.WeatherRequest{City: city})
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
