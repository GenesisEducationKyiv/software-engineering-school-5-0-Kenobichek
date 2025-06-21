package routes

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type httpRouteManager interface {
	http.Handler

	Route(pattern string, fn func(r chi.Router)) chi.Router
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Handle(pattern string, h http.Handler)
}

func NewHTTPRouter() chi.Router {
	return chi.NewRouter()
}

type weatherManager interface {
	GetWeather(writer http.ResponseWriter, request *http.Request)
}

type subscriptionManager interface {
	Subscribe(writer http.ResponseWriter, request *http.Request)
	Unsubscribe(writer http.ResponseWriter, request *http.Request)
	Confirm(writer http.ResponseWriter, request *http.Request)
}

type Service struct {
	router    httpRouteManager
	subscribe subscriptionManager
	weather   weatherManager
}

func NewService(
	weather weatherManager,
	subscribe subscriptionManager,
	router httpRouteManager,
) *Service {
	return &Service{
		router:    router,
		subscribe: subscribe,
		weather:   weather,
	}
}

func (r *Service) GetRouter() http.Handler {
	return r.router
}

func (r *Service) RegisterRoutes() {
	srv := r

	r.router.Route("/api", func(rt chi.Router) {
		rt.Get("/weather", srv.weather.GetWeather)
		rt.Post("/subscribe", srv.subscribe.Subscribe)
		rt.Get("/confirm/{token}", srv.subscribe.Confirm)
		rt.Get("/unsubscribe/{token}", srv.subscribe.Unsubscribe)
	})

	r.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.StripPrefix("/", http.FileServer(http.Dir("public")))
	r.router.Handle("/*", fs)
}
