package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

const (
	subscriptionServiceGroupID = "subscription-service"
)

type dbManagerImpl interface {
	InitDB(dsn string) error
	RunMigrations(migrationsPath string) error
	GetDB() *sql.DB
}

type loggerManager interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Sync() error
}

func Run(ctx context.Context, logger loggerManager) error {
	logger.Info("Subscription Service starting...")

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dbManager := infrastructure.NewDBManager(nil, logger)
	if err := dbManager.InitDB(cfg.GetDatabaseDSN()); err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer func() {
		if err := dbManager.GetDB().Close(); err != nil {
			logger.Error("db close error: %v", err)
		}
	}()

	if err := runMigrations(dbManager, "internal/migrations"); err != nil {
		return err
	}

	repo := subscriptions.New(dbManager.GetDB())
	publisher := infrastructure.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.Kafka.EventTopic)
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Error("publisher close error: %v", err)
		}
	}()

	dispatcher := handlers.NewDispatcher(repo, publisher, logger)

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
		subscriptionServiceGroupID,
		eventHandler,
		logger,
	)
	go consumer.Start(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	weatherClient, err := weatherclient.New(cfg.WeatherServiceAddr)
	if err != nil {
		logger.Error("failed to init weather client: %v", err)
		return fmt.Errorf("failed to init weather client: %w", err)
	}

	weatherJob := jobs.NewWeatherUpdateJob(repo, publisher, weatherClient, logger)
	go weatherJob.StartPeriodic(ctx)

	logger.Info("Subscription Service is running.")

	<-ctx.Done()
	logger.Info("Subscription Service shutting down...")

	return nil
}

func runMigrations(dbManager dbManagerImpl, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("migrations path %s does not exist", path)
	}
	if err := dbManager.RunMigrations(path); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
