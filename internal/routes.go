package internal

import (
	"net/http"

	"Weather-Forecast-API/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/weather", handlers.GetWeather)
		r.Post("/subscribe", handlers.Subscribe)
		r.Get("/confirm/{token}", handlers.Confirm)
		r.Get("/unsubscribe/{token}", handlers.Unsubscribe)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.StripPrefix("/", http.FileServer(http.Dir("public")))
	router.Handle("/*", fs)
}
