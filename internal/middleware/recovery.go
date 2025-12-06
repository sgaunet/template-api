package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sgaunet/template-api/internal/apperror"
)

var errPanic = errors.New("panic recovered")

// Recovery recovers from panics and returns structured error.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log stack trace
				debug.PrintStack()

				// Return internal error
				apperror.WriteError(w, apperror.NewInternalError(
					fmt.Errorf("%w: %v", errPanic, err),
				))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
