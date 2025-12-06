package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse is the HTTP error response format.
type ErrorResponse struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// WriteError writes a structured error response.
func WriteError(w http.ResponseWriter, err error) {
	var appErr *AppError
	var status int
	var response ErrorResponse

	// Convert to AppError if possible
	if errors.As(err, &appErr) {
		response = ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}
		status = errorCodeToHTTPStatus(appErr.Code)
	} else {
		// Unknown error - return 500
		response = ErrorResponse{
			Code:    ErrCodeInternal,
			Message: "An unexpected error occurred",
		}
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Failed to encode error response, nothing we can do
		return
	}
}

func errorCodeToHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeInternal:
		return http.StatusInternalServerError
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
