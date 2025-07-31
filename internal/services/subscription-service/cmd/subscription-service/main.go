package main

import (
	"context"
	"log"
	"subscription-service/internal/app"
	"subscription-service/internal/observability/logger"
)

func main() {
	ctx := context.Background()
	logger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("logger sync error: %v", err)
		}
	}()
	
	if err := app.Run(ctx, logger); err != nil {
		logger.Error("service exited with error: %v", err)
	}
}
