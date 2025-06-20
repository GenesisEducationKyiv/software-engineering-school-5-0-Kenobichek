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

func (a *App) ensureConfig() (Config, error) {
	if (a.config == Config{}) {
		return config.MustLoad()
	}
	return a.config, nil
}

func initDatabase(dsn string) (*sql.DB, error) {
	dbConn, err := db.Init(dsn)
	if err != nil {
		return nil, err
	}
	if err := db.RunMigrations(dbConn); err != nil {
		return nil, err
	}
	return dbConn, nil
}

func buildEmailNotifier(cfg Config) notifier.EmailNotifier {
	sgClient := sendgrid.NewSendClient(cfg.SendGrid.APIKey)
	sgNotifier := sendgridemailapi.NewSendgridNotifier(sgClient, cfg)
	return notifier.NewSendGridEmailNotifier(sgNotifier)
}

func buildWeatherProvider(cfg Config, client *http.Client) weatherprovider.WeatherProvider {
	geoSvc := openweather.NewOpenWeatherGeocodingService(cfg, client)
	owAPI := openweather.NewOpenWeatherAPI(cfg, client)
	return weatherprovider.NewOpenWeatherProvider(geoSvc, owAPI)
}

func buildHTTPRouter(
	weatherProv weatherprovider.WeatherProvider,
	subSvc subscription.SubscriptionService,
	notifSvc notification.NotificationService,
) http.Handler {
	rtr := routes.NewHTTPRouter()
	weatherHandler := weather.NewWeatherHandler(weatherProv, 5*time.Second)
	subscribeHandler := subscribe.NewSubscribeHandler(subSvc, notifSvc)
	appRouter := routes.NewRouter(weatherHandler, subscribeHandler, rtr)
	appRouter.RegisterRoutes()
	return appRouter.GetRouter()
}

func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func (a *App) startScheduler(s *scheduler.Scheduler, errCh chan<- error) {
	go func() {
		if _, err := s.Start(); err != nil {
			errCh <- err
		}
	}()
}
