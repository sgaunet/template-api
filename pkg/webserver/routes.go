package webserver

import (
	"net/http"
)

// "github.com/go-redis/redis/v7"

func (w *WebServer) initRoutes() {
	// Health check
	w.router.Get("/", HealthCheck)

	// Authors routes
	w.router.Post("/authors", w.authorsHandler.Create)
	w.router.Get("/authors", w.authorsHandler.List)
	w.router.Delete("/authors/{uuid}", w.authorsHandler.Delete)
}

// HealthCheck is the health check endpoint.
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
