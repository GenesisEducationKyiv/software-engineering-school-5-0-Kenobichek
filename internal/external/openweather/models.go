package weather

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type OpenWeatherData struct {
	Temperature float64
	Humidity    float64
	Description string
}

type OpenWeatherConfig struct {
	APIKey string
}

type WeatherData struct {
	Temperature float64
	Humidity    float64
	Description string
}
