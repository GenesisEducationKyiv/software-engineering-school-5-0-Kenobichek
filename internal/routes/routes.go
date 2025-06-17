package routes

import (
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"Weather-Forecast-API/internal/weather_provider"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Router struct {
	router       chi.Router
	subService   subscription.SubscriptionService
	notifService notification.NotificationService
	weatherProvider	weather_provider.WeatherProvider
}

func NewRouter(subService subscription.SubscriptionService, notifService notification.NotificationService) *Router {
	return &Router{
		router:       chi.NewRouter(),
		subService:   subService,
		notifService: notifService,
	}
}

func (r *Router) GetRouter() chi.Router {
	return r.router
}

func (r *Router) RegisterRoutes() {
	subscribeHandler := subscribe.NewSubscribeHandler(
		r.subService,
		r.notifService,
	)
	weatherHandler := weather.NewWeatherHandler(r.weatherProvider, 5 * time.Second )

	r.router.Route("/api", func(r chi.Router) {
		r.Get("/weather", weatherHandler.GetWeather)
		r.Post("/subscribe", subscribeHandler.Subscribe)
		r.Get("/confirm/{token}", subscribeHandler.Confirm)
		r.Get("/unsubscribe/{token}", subscribeHandler.Unsubscribe)
	})

	r.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.StripPrefix("/", http.FileServer(http.Dir("public")))
	r.router.Handle("/*", fs)
}
