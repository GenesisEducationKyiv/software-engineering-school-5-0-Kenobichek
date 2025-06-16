package routes

import (
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterRoutes(router chi.Router) {
	subscribeHandler := subscribe.NewSubscribeHandler()
	weatherHandler := weather.NewWeatherHandlerWithDefault()

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
