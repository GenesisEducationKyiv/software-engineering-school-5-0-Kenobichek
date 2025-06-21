package app

import (
	"Weather-Forecast-API/config"
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/external/sendgridemailapi"
	"Weather-Forecast-API/internal/db"
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/scheduler"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"Weather-Forecast-API/internal/weatherprovider"
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

func (a *App) buildEmailNotifier() notifier.EmailNotifier {
	sgCfg := a.config.SendGrid

	sgClient := sendgrid.NewSendClient(sgCfg.APIKey)

	sgNotifier := sendgridemailapi.NewSendgridNotifier(
		sgClient,
		sgCfg.SenderName,
		sgCfg.SenderEmail,
	)

	return notifier.NewSendGridEmailNotifier(sgNotifier)
}

func (a *App) buildWeatherProvider(client *http.Client) weatherprovider.WeatherProvider {
	owCfg := a.config.OpenWeather

	geoSvc := openweather.NewOpenWeatherGeocodingService(
		client,
		owCfg.GeocodingAPIURL,
		owCfg.APIKey,
	)

	owAPI := openweather.NewOpenWeatherAPI(
		client,
		owCfg.WeatherAPIURL,
		owCfg.APIKey,
	)
	return weatherprovider.NewOpenWeatherProvider(geoSvc, owAPI)
}

func (a *App) buildHTTPRouter(
	weatherProv weatherprovider.WeatherProvider,
	subSvc subscription.SubscriptionService,
	notifSvc notification.NotificationService,
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

	appRouter := routes.NewRouter(
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
