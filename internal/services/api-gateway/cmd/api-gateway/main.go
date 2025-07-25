package main

import (
	"api-gateway/app"
	"log"
)

func main() {
	log.Println("Starting API gateway…")
	if err := app.Run(); err != nil {
		log.Fatalf("api-gateway exited with error: %v", err)
	}
	log.Println("API gateway stopped")
}
