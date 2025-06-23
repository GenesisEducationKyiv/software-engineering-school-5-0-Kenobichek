package app

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/external/sendgridemailapi"
	"Weather-Forecast-API/internal/db"
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/notifier/sengridnotifier"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/scheduler"
	"Weather-Forecast-API/internal/weatherprovider/openweatherprovider"
	"context"
	"database/sql"
	"github.com/sendgrid/sendgrid-go"
	"net/http"
	"time"
)

const (
	readTimeout           = 5 * time.Second
	writeTimeout          = 10 * time.Second
	idleTimeout           = 120 * time.Second
	weatherHandlerTimeout = 5 * time.Second
)

func (a *App) loadConfigIfEmpty() (Config, error) {
	if (a.config == Config{}) {
		return config.MustLoad()
	}

	return a.config, nil
}

func (a *App) connectDatabase() (*sql.DB, error) {
	dbConn, err := db.Init(a.config.GetDatabaseDSN())
	if err != nil {
		return nil, err
	}

	if err := db.RunMigrations(dbConn); err != nil {
		return nil, err
	}
	return dbConn, nil
}

func (a *App) buildSendGridEmailNotifier() *sengridnotifier.SendGridEmailNotifier {
	sgCfg := a.config.SendGrid

	sgClient := sendgrid.NewSendClient(sgCfg.APIKey)

	sgNotifier := sendgridemailapi.NewSendgridNotifier(
		sgClient,
		sgCfg.SenderName,
		sgCfg.SenderEmail,
	)

	return sengridnotifier.NewSendGridEmailNotifier(sgNotifier)
}

func (a *App) buildOpenWeatherProvider(client *http.Client) *openweatherprovider.OpenWeatherProvider {
	owCfg := a.config.OpenWeather

	geoSvc := openweather.NewGeocodingService(
		client,
		owCfg.GeocodingAPIURL,
		owCfg.APIKey,
	)

	owAPI := openweather.NewOpenWeatherAPI(
		client,
		owCfg.WeatherAPIURL,
		owCfg.APIKey,
	)
	return openweatherprovider.NewOpenWeatherProvider(geoSvc, owAPI)
}

type subscriptionManager interface {
	Subscribe(sub *subscriptions.Info) error
	Unsubscribe(sub *subscriptions.Info) error
	Confirm(sub *subscriptions.Info) error
	GetSubscriptionByToken(token string) (*subscriptions.Info, error)
	GetDueSubscriptions() []subscriptions.Info
	UpdateNextNotification(id int, next time.Time) error
}

type notificationManager interface {
	SendWeatherUpdate(channel string, recipient string, metrics weather.Metrics) error
	SendConfirmation(channel string, recipient string, token string) error
	SendUnsubscribe(channel string, recipient string, city string) error
}

type weatherProviderManager interface {
	GetWeatherByCity(ctx context.Context, city string) (weather.Metrics, error)
}

func (a *App) buildRouter(
	weatherProv weatherProviderManager,
	subSvc subscriptionManager,
	notifSvc notificationManager,
) http.Handler {
	rtr := routes.NewHTTPRouter()

	weatherHandler := weather.NewHandler(
		weatherProv,
		weatherHandlerTimeout,
	)

	subscribeHandler := subscribe.NewHandler(
		subSvc,
		notifSvc,
	)

	appRouter := routes.NewService(
		weatherHandler,
		subscribeHandler,
		rtr,
	)

	appRouter.RegisterRoutes()

	return appRouter.GetRouter()
}

func (a *App) newHTTPServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         a.config.GetServerAddress(),
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func (a *App) runSchedulerAsync(s *scheduler.Scheduler, errCh chan<- error) {
	go func() {
		if _, err := s.Start(); err != nil {
			errCh <- err
		}
	}()
}
