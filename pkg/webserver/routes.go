package webserver

import (
	"net/http"

	"github.com/sgaunet/template-api/api/authors"
)

// "github.com/go-redis/redis/v7"

func (w *WebServer) initRoutes(authorSvc *authors.Service) {
	w.router.Post("/authors", authorSvc.Create)
	w.router.Get("/authors/{id}", authorSvc.Get)
	w.router.Put("/authors/{id}", authorSvc.FullUpdate)
	w.router.Delete("/authors/{id}", authorSvc.Delete)
	w.router.Get("/authors", authorSvc.List)
	w.router.Get("/", HealthCheck)
}

// HealthCheck is the health check endpoint
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
