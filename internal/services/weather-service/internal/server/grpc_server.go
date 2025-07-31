package server

import (
	"context"
	"fmt"
	"net"

	"internal/services/weather-service/internal/provider"
	"internal/services/weather-service/proto"

	"google.golang.org/grpc"
)

type WeatherGRPCServer struct {
	proto.UnimplementedWeatherServiceServer
	provider *provider.CachedWeatherProvider
}

func NewWeatherGRPCServer(provider *provider.CachedWeatherProvider) *WeatherGRPCServer {
	return &WeatherGRPCServer{provider: provider}
}

func (s *WeatherGRPCServer) GetWeather(ctx context.Context, req *proto.WeatherRequest) (*proto.WeatherResponse, error) {
	if req.GetCity() == "" {
		return nil, fmt.Errorf("city is required")
	}
	metrics, err := s.provider.GetWeatherByCity(ctx, req.GetCity())
	if err != nil {
		return nil, err
	}
	return &proto.WeatherResponse{
		City:        metrics.City,
		Description: metrics.Description,
		Temperature: metrics.Temperature,
		Humidity:    metrics.Humidity,
	}, nil
}

func RunGRPCServer(address string, provider *provider.CachedWeatherProvider) error {
	var lc net.ListenConfig
	lis, err := lc.Listen(context.Background(), "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterWeatherServiceServer(grpcServer, NewWeatherGRPCServer(provider))
	return grpcServer.Serve(lis)
}
