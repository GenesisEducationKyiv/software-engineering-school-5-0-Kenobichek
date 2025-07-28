package weatherclient

import (
	"context"
	"fmt"
	"log"

	"subscription-service/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type connectionManager interface {
	Close() error
}

type WeatherClient struct {
    client proto.WeatherServiceClient
	conn connectionManager
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
	return &WeatherClient{
		client: proto.NewWeatherServiceClient(conn),
		conn:   conn,
	}, nil
}	
	
func (a *WeatherClient) GetWeather(ctx context.Context, req *proto.WeatherRequest) (*proto.WeatherResponse, error) {
	return a.client.GetWeather(ctx, req)
}

func (w *WeatherClient) Close() error {
	return w.conn.Close()
}
	