package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"subscription-service/internal/observability/metrics"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	metricsPath = "/metrics"
	metricsReadTimeout = 5 * time.Second
	metricsWriteTimeout = 10 * time.Second
	metricsIdleTimeout = 120 * time.Second
)

type dbManagerImpl interface {
	InitDB(dsn string) error
	RunMigrations(migrationsPath string) error
	GetDB() *sql.DB
}

type loggerManager interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Sync() error
}

func Run(ctx context.Context, logger loggerManager) error {
	logger.Infof("Subscription Service starting...")

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := metrics.Register(); err != nil {
		logger.Errorf("failed to register metrics", "error", err)
	}

	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Observability.VictoriaMetricsPort),
		Handler:      mux,
		ReadTimeout:  metricsReadTimeout,
		WriteTimeout: metricsWriteTimeout,
		IdleTimeout:  metricsIdleTimeout,
	}
	
	metricsServerErr := make(chan error, 1)
	go func() {
		logger.Infof("metrics endpoint listening", "addr", fmt.Sprintf(":%d", cfg.Observability.VictoriaMetricsPort))
		if err := metricsServer.ListenAndServe(); 
			err != nil && err != http.ErrServerClosed {
			metricsServerErr <- fmt.Errorf("metrics server error: %w", err)
		}
	}()

	dbManager := infrastructure.NewDBManager(nil, logger)
	if err := dbManager.InitDB(cfg.GetDatabaseDSN()); err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer func() {
		if err := dbManager.GetDB().Close(); err != nil {
			logger.Errorf("db close error: %v", err)
		}
	}()

	if err := runMigrations(dbManager, "internal/migrations"); err != nil {
		return err
	}

	repo := subscriptions.New(dbManager.GetDB())
	publisher := infrastructure.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.Kafka.EventTopic)
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Errorf("publisher close error: %v", err)
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
	consumerDone := consumer.Start(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	weatherClient, err := weatherclient.New(cfg.WeatherServiceAddr)
	if err != nil {
		logger.Errorf("failed to init weather client: %v", err)
		return fmt.Errorf("failed to init weather client: %w", err)
	}

	weatherJob := jobs.NewWeatherUpdateJob(repo, publisher, weatherClient, logger)
	go weatherJob.StartPeriodic(ctx)

	logger.Infof("Subscription Service is running.")

	<-ctx.Done()
	logger.Infof("Subscription Service shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Infof("Shutting down metrics server...")
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("metrics server shutdown error", "error", err)
	}
	select {
	case err := <-metricsServerErr:
		logger.Errorf("metrics server error", "error", err)
	default:
	}

	logger.Infof("Waiting for Kafka consumer to finish...")
	select {
	case <-consumerDone:
		logger.Infof("Kafka consumer stopped")
	case <-shutdownCtx.Done():
		logger.Errorf("Kafka consumer shutdown timeout", "error", shutdownCtx.Err())
	}

	logger.Infof("Weather job stopped (by context)")

	logger.Infof("Kafka publisher closed (deferred)")

	logger.Infof("DB connection closed (deferred)")

	logger.Infof("Flushing logger...")
	if err := logger.Sync(); err != nil {
		logger.Errorf("logger sync error", "error", err)
	}

	logger.Infof("Graceful shutdown completed")
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
