package main

import (
	"api-gateway/app"
	"log"
)

func main() {
	log.Println("api-gateway starting...")
	if err := app.Run(); err != nil {
		log.Fatalf("api-gateway exited with error: %v", err)
	}
	log.Println("api-gateway stopped")
}
