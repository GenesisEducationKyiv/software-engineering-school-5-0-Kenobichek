package weatherclient

import (
	"context"
	"fmt"
	"log"

	"subscription-service/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WeatherClient struct {
    proto.WeatherServiceClient
}

func New(addr string) (*WeatherClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial weather service at %s: %w", addr, err)
	}
	log.Printf("[WeatherHandler] gRPC connection established to %s", addr)
	return &WeatherClient{WeatherServiceClient: proto.NewWeatherServiceClient(conn)}, nil
}	
	
func (a *WeatherClient) GetWeather(ctx context.Context, req *proto.WeatherRequest) (*proto.WeatherResponse, error) {
	return a.WeatherServiceClient.GetWeather(ctx, req)
}
