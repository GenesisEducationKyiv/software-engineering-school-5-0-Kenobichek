package weatherclient

import (
	"context"
	"fmt"
	"log"

	"api-gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WeatherClient struct {
	client proto.WeatherServiceClient
	conn   *grpc.ClientConn
}

func New(addr string) (*WeatherClient, error) {
	if addr == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial weather service at %s: %w", addr, err)
	}
	log.Printf("[WeatherClient] gRPC connection established to %s", addr)
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
