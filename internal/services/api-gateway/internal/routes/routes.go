package routes

import (
	"api-gateway/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, weatherHandler *handlers.WeatherHandler, subscribeHandler *handlers.SubscribeHandler) {
	r.Get("/weather", weatherHandler.WeatherProxyHandler)

	r.Post("/subscribe", subscribeHandler.Subscribe)
	r.Post("/confirm", subscribeHandler.ConfirmSubscription)
	r.Post("/unsubscribe", subscribeHandler.Unsubscribe)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
}
