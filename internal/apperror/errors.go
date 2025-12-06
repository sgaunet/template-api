// Package apperror provides structured error handling for the application.
package apperror

import (
	"errors"
	"fmt"
)

// ErrorCode represents application error codes.
type ErrorCode string

const (
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
)

// AppError represents a structured application error.
type AppError struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error.
func NewValidationError(message string, details map[string]string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: message,
	}
}

// NewInternalError creates a new internal error.
func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: "An internal error occurred",
		Err:     err,
	}
}

// NewBadRequestError creates a new bad request error.
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
	}
}

// NewConflictError creates a new conflict error.
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
	}
}

// IsValidationError checks if error is a validation error.
func IsValidationError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeValidation
	}
	return false
}

// IsNotFoundError checks if error is a not found error.
func IsNotFoundError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}
