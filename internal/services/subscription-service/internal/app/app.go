package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscription-service/config"
	"subscription-service/internal/handlers"
	"subscription-service/internal/infrastructure"
	"subscription-service/internal/repository"
	"subscription-service/internal/jobs"
	"subscription-service/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	migrationsPath := "internal/migrations"
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

	grpcAddr := cfg.WeatherServiceAddr

	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("[WeatherHandler] failed to dial gRPC at %s: %v", grpcAddr, err)
		return err
	}
	log.Printf("[WeatherHandler] gRPC connection established to %s", grpcAddr)
	weatherClient := proto.NewWeatherServiceClient(conn)

	weatherJob := jobs.NewWeatherUpdateJob(repo, publisher, weatherClient)
	go weatherJob.StartPeriodic(ctx)

	log.Println("Subscription Service is running.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Subscription Service shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
	return nil
}
