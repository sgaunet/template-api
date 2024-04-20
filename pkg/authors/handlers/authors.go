package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sgaunet/template-api/domain"
	"github.com/sgaunet/template-api/internal/httperror"
	authorsservice "github.com/sgaunet/template-api/pkg/authors/service"
)

// AuthorHandlers is the Authors Handlers
type AuthorHandlers struct {
	authorService *authorsservice.AuthorService
}

// NewAuthorsHandlers creates a new authors handlers
func NewAuthorsHandlers(authrSvc *authorsservice.AuthorService) *AuthorHandlers {
	return &AuthorHandlers{
		authorService: authrSvc,
	}
}

// Create creates a new author
func (a *AuthorHandlers) Create(w http.ResponseWriter, r *http.Request) {
	// Parse body request
	payload := domain.AuthorRequestBody{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &payload)
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}
	// Validate
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(payload)
	// validationErrors := err.(validator.ValidationErrors)
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}

	author, err := a.authorService.Create(context.Background(), payload.Name, payload.Bio)
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}
	response := domain.FromDB(author)
	w.WriteHeader(http.StatusCreated)
	// write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		httperror.WriteStatusInternalServerError(w, err)
		return
	}
}

// List returns all authors
func (a *AuthorHandlers) List(w http.ResponseWriter, r *http.Request) {
	authors, err := a.authorService.List(context.Background())
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}
	// write response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		httperror.WriteStatusInternalServerError(w, err)
	}
}

// Delete deletes an author by id
func (a *AuthorHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	uuidParam := chi.URLParam(r, "uuid")
	// Convert to int64
	id, err := strconv.ParseInt(uuidParam, 10, 64)
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}
	err = a.authorService.Delete(r.Context(), id)
	if err != nil {
		httperror.WriteBadRequestError(w, err)
		return
	}
}
