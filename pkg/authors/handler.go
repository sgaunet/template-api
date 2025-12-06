// Package authors provides the authors domain logic.
package authors

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sgaunet/template-api/internal/apperror"
)

// Handler handles HTTP requests for authors.
type Handler struct {
	service Service
}

// NewHandler creates a new author handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /authors.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAuthorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.NewBadRequestError("Invalid request body"))
		return
	}
	defer r.Body.Close()

	author, err := h.service.Create(r.Context(), &req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(author.ToResponse())
}

// List handles GET /authors.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.List(r.Context())
	if err != nil {
		apperror.WriteError(w, err)
		return
	}

	responses := make([]*AuthorResponse, len(authors))
	for i, author := range authors {
		responses[i] = author.ToResponse()
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(responses)
}

// Delete handles DELETE /authors/{uuid}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "uuid")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		apperror.WriteError(w, apperror.NewBadRequestError("Invalid author ID"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		apperror.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
