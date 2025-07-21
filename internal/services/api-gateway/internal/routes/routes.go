package routes

import (
	"api-gateway/internal/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	weatherHandler *handlers.WeatherHandler,
	subscribeHandler *handlers.SubscribeHandler,
) {
	r.Get("/weather", weatherHandler.WeatherProxyHandler)

	r.Post("/subscribe", subscribeHandler.Subscribe)
	r.Get("/confirm/{token}", subscribeHandler.ConfirmSubscription)
	r.Get("/unsubscribe/{token}", subscribeHandler.Unsubscribe)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	})
}
