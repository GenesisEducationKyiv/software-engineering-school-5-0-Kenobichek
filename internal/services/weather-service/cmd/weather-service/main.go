package main

import (
	"internal/services/weather-service/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("service exited with error: %v", err)
	}
}
