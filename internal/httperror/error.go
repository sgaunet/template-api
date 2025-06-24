// Package httperror provides helper functions for handling HTTP errors.
package httperror

import "net/http"

// WriteStatusInternalServerError writes a 500 error and message.
func WriteStatusInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

// WriteBadRequestError writes a bad request status code and the error message to the response writer.
func WriteBadRequestError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
}
