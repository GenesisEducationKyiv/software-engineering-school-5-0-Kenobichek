package app

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/internal/httpclient"
	"Weather-Forecast-API/internal/scheduler"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"errors"
	"log"
	"net/http"
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
	dbConn, err := a.connectDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dbConn.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
	}()

	emailNotifier := a.buildEmailNotifier()
	httpClient := httpclient.New()
	weatherProv := a.buildWeatherProvider(httpClient)

	subSvc := subscription.NewSubscriptionService()
	notifSvc := notification.NewNotificationService(emailNotifier)

	taskScheduler := scheduler.NewScheduler(notifSvc, weatherProv)
	errCh := make(chan error, 1)

	a.runSchedulerAsync(taskScheduler, errCh)

	router := a.buildHTTPRouter(weatherProv, subSvc, notifSvc)

	server := a.newHTTPServer(router)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	return <-errCh
}
