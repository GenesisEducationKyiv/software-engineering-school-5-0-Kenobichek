package weather

type DataWeather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

type Provider interface {
	GetWeather(city string) (DataWeather, error)
}
