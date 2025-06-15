package weather

import "context"

type DataWeather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

type Provider interface {
	GetWeather(ctx context.Context, city string) (DataWeather, error)
}
