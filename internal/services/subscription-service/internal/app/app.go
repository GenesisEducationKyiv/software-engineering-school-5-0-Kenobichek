package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"subscription-service/internal/observability"
	"os"
	"os/signal"
	"syscall"

	"subscription-service/config"
	"subscription-service/internal/domain"
	"subscription-service/internal/handlers"
	"subscription-service/internal/infrastructure"
	"subscription-service/internal/jobs"
	"subscription-service/internal/repository/subscriptions"
	"subscription-service/internal/weatherclient"
)

func Run(ctx context.Context) error {
	observability.Infof("Subscription Service starting...")

	cfg, err := config.MustLoad()
	if err != nil {	
		return fmt.Errorf("load config: %w", err)
	}

	// initialize metrics server early
	observability.StartMetricsServer()

	db, err := infrastructure.InitDB(cfg.GetDatabaseDSN())
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			observability.Warnf("db close error: %v", err)
		}
	}()

	if err := runMigrations(db, "internal/migrations"); err != nil {
		return err
	}

	repo := subscriptions.New(db)
	publisher := infrastructure.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.Kafka.EventTopic)
	defer func() {
		if err := publisher.Close(); err != nil {
			observability.Warnf("publisher close error: %v", err)
		}
	}()

	dispatcher := handlers.NewDispatcher(repo, publisher)

	eventHandler := func(ctx context.Context, topic string, message []byte) error {
		var cmd domain.SubscriptionCommand
		if err := json.Unmarshal(message, &cmd); err != nil {
			return fmt.Errorf("unmarshal subscription command: %w", err)
		}
		return dispatcher.Handle(ctx, cmd)
	}
	topics := []string{cfg.Kafka.CommandTopic}

	consumer := infrastructure.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		topics,
		"subscription-service",
		eventHandler,
	)
	go consumer.Start(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	weatherClient, err := weatherclient.New(cfg.WeatherServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to init weather client: %w", err)
	}

	weatherJob := jobs.NewWeatherUpdateJob(repo, publisher, weatherClient)
	go weatherJob.StartPeriodic(ctx)

	observability.Infof("Subscription Service is running.")

	<-ctx.Done()
	observability.Infof("Subscription Service shutting down...")

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
