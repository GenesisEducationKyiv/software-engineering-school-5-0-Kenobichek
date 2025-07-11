package weather

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Metrics struct {
	Temperature float64
	Humidity    float64
	Description string
	City        string
}
