package webserver

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sgaunet/template-api/api/authors"
	// "github.com/go-redis/redis/v7"
)

const listenAddr string = ":3000"

// WebServer is the web server
type WebServer struct {
	srv    *http.Server
	router *chi.Mux
}

// NewWebServer creates a new web server
func NewWebServer(authorsSvc *authors.Service) (*WebServer, error) {
	w := &WebServer{}
	w.router = chi.NewRouter()
	w.router.Use(middleware.RequestID)
	w.router.Use(middleware.Logger)
	w.router.Use(middleware.Recoverer)
	w.initRoutes(authorsSvc)

	w.srv = &http.Server{
		Addr:              listenAddr,
		Handler:           w.router,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}
	return w, nil
}

// Start shuts down the web server
func (w *WebServer) Start() error {
	err := w.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown shuts down the web server
func (w *WebServer) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

// SetListenAddr sets the listen address (format expected: ":3000")
// It won't restart the webserver if it's already running
func (w *WebServer) SetListenAddr(addr string) {
	w.srv.Addr = addr
}
