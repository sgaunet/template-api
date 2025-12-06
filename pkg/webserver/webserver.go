// Package webserver provides the web server for the application.
package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sgaunet/template-api/internal/middleware"
	"github.com/sgaunet/template-api/pkg/authors"
	// "github.com/go-redis/redis/v7".
)

const (
	listenAddr        string = ":3000"
	readHeaderTimeout        = 100 * time.Millisecond
)

// WebServer is the web server.
type WebServer struct {
	srv            *http.Server
	router         *chi.Mux
	authorsHandler *authors.Handler
}

// NewWebServer creates a new web server.
func NewWebServer(authorsHandler *authors.Handler) (*WebServer, error) {
	w := &WebServer{
		authorsHandler: authorsHandler,
	}
	w.router = chi.NewRouter()

	// Global middleware
	w.router.Use(chimiddleware.RequestID)
	w.router.Use(chimiddleware.Logger)
	w.router.Use(middleware.Recovery)
	w.router.Use(middleware.JSONContentType)

	w.initRoutes()

	w.srv = &http.Server{
		Addr:              listenAddr,
		Handler:           w.router,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	return w, nil
}

// Start starts the web server.
func (w *WebServer) Start() error {
	err := w.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("could not start webserver: %w", err)
	}
	return nil
}

// Shutdown shuts down the web server.
func (w *WebServer) Shutdown(ctx context.Context) error {
	if err := w.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("could not shutdown webserver: %w", err)
	}
	return nil
}

// SetListenAddr sets the listen address (format expected: ":3000")
// It won't restart the webserver if it's already running.
func (w *WebServer) SetListenAddr(addr string) {
	w.srv.Addr = addr
}
