package app

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/httpclient"
	"Weather-Forecast-API/internal/repository/emailtemplates"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/scheduler"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config = config.Config

type App struct {
	config Config
}

func New(cfg Config) *App { return &App{config: cfg} }

func (a *App) Run() error {
	var err error
	a.config, err = a.loadConfigIfEmpty()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn, err := a.connectDatabase()
	if err != nil {
		return err
	}

	emailNotifier := a.buildSendGridEmailNotifier()
	httpClient := httpclient.New()

	weatherProvChain, err := a.buildWeatherProviderChain(httpClient)
	if err != nil {
		return err
	}

	subsRepo := subscriptions.New(dbConn)
	tmplsRepo := emailtemplates.New(dbConn)

	subSvc := subscription.NewService(subsRepo)
	notifSvc := notification.NewService(emailNotifier, tmplsRepo)

	taskScheduler := scheduler.NewScheduler(notifSvc, subSvc, weatherProvChain, weatherHandlerTimeout)

	router := a.buildRouter(weatherProvChain, subSvc, notifSvc)
	server := a.newHTTPServer(router)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := taskScheduler.Start(); err != nil {
			log.Printf("Scheduler error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Printf("Server is running on %s", a.config.GetServerAddress())

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := taskScheduler.Stop(); err != nil {
		log.Printf("Scheduler shutdown error: %v", err)
	}

	wg.Wait()

	if err := dbConn.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}

	log.Println("Application shutdown complete")
	return nil
}
