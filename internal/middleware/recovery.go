// Package middleware provides HTTP middleware for the application.
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sgaunet/template-api/internal/apperror"
)

// Recovery recovers from panics and returns structured error.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log stack trace
				debug.PrintStack()

				// Return internal error
				apperror.WriteError(w, apperror.NewInternalError(
					fmt.Errorf("panic: %v", err),
				))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
