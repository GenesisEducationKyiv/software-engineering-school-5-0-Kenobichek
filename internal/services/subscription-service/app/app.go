package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscription-service/config"
	"subscription-service/handlers"
	"subscription-service/infrastructure"
	"subscription-service/repository"
)

func Run(ctx context.Context) error {
	log.Println("Subscription Service starting...")

	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	db, err := infrastructure.InitDB(cfg.GetDatabaseDSN())
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("DB close error: %v", err)
		}
	}()

	migrationsPath := "migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return err
	}
	if err := infrastructure.RunMigrations(db, migrationsPath); err != nil {
		return err
	}

	repo := repository.New(db)
	publisher := infrastructure.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.Kafka.EventTopic)
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("Publisher close error: %v", err)
		}
	}()
	dispatcher := handlers.NewDispatcher(repo, publisher)

	consumer := infrastructure.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		"commands.subscription",
		dispatcher.Handle,
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go consumer.Start(ctx)

	log.Println("Subscription Service is running.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Subscription Service shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
	return nil
}
