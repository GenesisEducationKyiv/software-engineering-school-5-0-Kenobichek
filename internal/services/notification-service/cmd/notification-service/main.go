package main

import (
	"context"
	"log"
	"notification-service/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		log.Fatalf("service exited with error: %v", err)
	}
}
