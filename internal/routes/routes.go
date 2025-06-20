package routes

import (
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type HTTPRouter interface {
	http.Handler

	Route(pattern string, fn func(r chi.Router)) chi.Router
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Handle(pattern string, h http.Handler)
}

func NewHTTPRouter() HTTPRouter {
	return chi.NewRouter()
}

type RouterManager interface {
	GetRouter() HTTPRouter
	RegisterRoutes()
}

type ServerRouter struct {
	router    HTTPRouter
	subscribe subscribe.SubscriptionManager
	weather   weather.WeatherManager
}

func NewRouter(
	weather weather.WeatherManager,
	subscribe subscribe.SubscriptionManager,
	router HTTPRouter) *ServerRouter {
	return &ServerRouter{
		router:    router,
		subscribe: subscribe,
		weather:   weather,
	}
}

func (r *ServerRouter) GetRouter() HTTPRouter {
	return r.router
}

func (r *ServerRouter) RegisterRoutes() {
	outer := r

	r.router.Route("/api", func(rt chi.Router) {
		rt.Get("/weather", outer.weather.GetWeather)
		rt.Post("/subscribe", outer.subscribe.Subscribe)
		rt.Get("/confirm/{token}", outer.subscribe.Confirm)
		rt.Get("/unsubscribe/{token}", outer.subscribe.Unsubscribe)
	})

	r.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.StripPrefix("/", http.FileServer(http.Dir("public")))
	r.router.Handle("/*", fs)
}
