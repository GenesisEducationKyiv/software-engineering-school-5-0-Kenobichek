package internal

import (
	"Weather-Forecast-API/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Get("/weather", handlers.GetWeather)
		r.Post("/subscribe", handlers.Subscribe)
		r.Get("/confirm/{token}", handlers.Confirm)
		r.Get("/unsubscribe/{token}", handlers.Unsubscribe)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.StripPrefix("/", http.FileServer(http.Dir("public")))
	r.Handle("/*", fs)

}
