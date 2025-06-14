package internal

import (
	weather "Weather-Forecast-API/internal/external/openweather"
	notificationService "Weather-Forecast-API/internal/services/notification"
	subscriptionService "Weather-Forecast-API/internal/services/subscription"
	"net/http"
	"os"

	"Weather-Forecast-API/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router) {
	subscribeHandler := handlers.NewSubscribeHandler(
		subscriptionService.NewSubscriptionService(),
		notificationService.NewNotificationService())

	weatherHandler := handlers.NewWeatherHandler(
		*weather.NewOpenWeatherProvider(os.Getenv("OPENWEATHERMAP_API_KEY")))

	router.Route("/api", func(r chi.Router) {
		r.Get("/weather", weatherHandler.GetWeather)
		r.Post("/subscribe", subscribeHandler.Subscribe)
		r.Get("/confirm/{token}", subscribeHandler.Confirm)
		r.Get("/unsubscribe/{token}", subscribeHandler.Unsubscribe)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.StripPrefix("/", http.FileServer(http.Dir("public")))
	router.Handle("/*", fs)
}
