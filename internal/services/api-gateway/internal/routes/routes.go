package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type weatherHandlerManager interface {
	WeatherProxyHandler(w http.ResponseWriter, r *http.Request)
}

type subscribeHandlerManager interface {
	Subscribe(w http.ResponseWriter, r *http.Request)
	ConfirmSubscription(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(
	r chi.Router,
	weatherHandler weatherHandlerManager,
	subscribeHandler subscribeHandlerManager,
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
