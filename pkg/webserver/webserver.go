// Package webserver provides the web server for the application.
package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sgaunet/template-api/pkg/authors/handlers"
	authors "github.com/sgaunet/template-api/pkg/authors/service"
	// "github.com/go-redis/redis/v7".
)

const (
	listenAddr        string = ":3000"
	readHeaderTimeout        = 100 * time.Millisecond
)

// WebServer is the web server.
type WebServer struct {
	srv      *http.Server
	router   *chi.Mux
	handlers *handlers.AuthorHandlers
}

// NewWebServer creates a new web server.
func NewWebServer(authorsSvc *authors.AuthorService) (*WebServer, error) {
	w := &WebServer{
		handlers: handlers.NewAuthorsHandlers(authorsSvc),
	}
	w.router = chi.NewRouter()
	w.router.Use(middleware.RequestID)
	w.router.Use(middleware.Logger)
	w.router.Use(middleware.Recoverer)
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
