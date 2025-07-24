package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"subscription-service/config"
	"subscription-service/internal/handlers"
	"subscription-service/internal/infrastructure"
	"subscription-service/internal/jobs"
	"subscription-service/internal/proto"
	"subscription-service/internal/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func Run(ctx context.Context) error {
	log.Println("Subscription Service starting...")

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := infrastructure.InitDB(cfg.GetDatabaseDSN())
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	if err := runMigrations(db, "internal/migrations"); err != nil {
		return err
	}

	repo := repository.New(db)
	publisher := infrastructure.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.Kafka.EventTopic)
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("publisher close error: %v", err)
		}
	}()

	dispatcher := handlers.NewDispatcher(repo, publisher)

	consumer := infrastructure.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.CommandTopic,
		dispatcher.Handle,
	)
	go consumer.Start(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	weatherClient, err := newWeatherClient(cfg.WeatherServiceAddr)
	if err != nil {
		return err
	}

	weatherJob := jobs.NewWeatherUpdateJob(repo, publisher, weatherClient)
	go weatherJob.StartPeriodic(ctx)

	log.Println("Subscription Service is running.")

	<-ctx.Done()
	log.Println("Subscription Service shutting down...")

	return nil
}

func runMigrations(db *sql.DB, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("migrations path %s does not exist", path)
	}
	if err := infrastructure.RunMigrations(db, path); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

func newWeatherClient(addr string) (proto.WeatherServiceClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("dial weather service at %s: %w", addr, err)
	}
	log.Printf("[WeatherHandler] gRPC connection established to %s", addr)
	return proto.NewWeatherServiceClient(conn), nil
}
