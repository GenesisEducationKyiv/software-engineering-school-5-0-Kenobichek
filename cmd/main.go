package main

import (
	"os"

	"Weather-Forecast-API/internal/app"
)

func main() {
	cfg := app.Config{}
	if err := app.New(cfg).Run(); err != nil {
		os.Exit(1)
	}
}
