package webserver

import (
	"net/http"
)

// "github.com/go-redis/redis/v7"

func (w *WebServer) initRoutes() {
	w.router.Post("/authors", w.handlers.Create)
	// w.router.Get("/authors/{id}", authorSvc.Get)
	// w.router.Put("/authors/{id}", authorSvc.FullUpdate)
	// w.router.Delete("/authors/{id}", authorSvc.Delete)
	w.router.Get("/authors", w.handlers.List)
	w.router.Get("/", HealthCheck)
}

// HealthCheck is the health check endpoint
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
