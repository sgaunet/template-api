package authors

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sgaunet/template-api/internal/database"
)

// Service is the Authors Service Layer
type Service struct {
	queries database.Querier
}

// NewService creates a new authors service
func NewService(queries database.Querier) *Service {
	return &Service{queries: queries}
}

// apiAuthor is the author representation for the API
type apiAuthor struct {
	ID   int64
	Name string `json:"name,omitempty" validate:"min=5,max=20"`
	Bio  string `json:"bio,omitempty"`
}

// Create creates a new author
func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	// Parse body request
	payload := apiAuthor{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// Validate
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(payload)
	// validationErrors := err.(validator.ValidationErrors)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	author, err := s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: payload.Name,
		Bio:  payload.Bio,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	response := fromDB(author)
	w.WriteHeader(http.StatusCreated)
	// write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func fromDB(author database.Author) *apiAuthor {
	return &apiAuthor{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}

// Get returns an author by id
func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("id is required"))
		return
	}
	// convert id to int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// Get author
	author, err := s.queries.GetAuthor(context.Background(), idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("id does not exist"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	response := fromDB(author)
	// application/json
	w.Header().Set("Content-Type", "application/json")
	// write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

// Delete deletes an author by id
func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("id is required"))
		return
	}
	// convert id to int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// Delete author
	if err := s.queries.DeleteAuthor(context.Background(), idInt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("id does not exist"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

// List returns all authors
func (s *Service) List(w http.ResponseWriter, _ *http.Request) {
	authors, err := s.queries.ListAuthors(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	response := []*apiAuthor{}
	for _, author := range authors {
		response = append(response, fromDB(author))
	}
	// write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

// FullUpdate updates an author by id
func (s *Service) FullUpdate(w http.ResponseWriter, r *http.Request) {
	payload := database.UpdateAuthorParams{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	author, err := s.queries.UpdateAuthor(context.Background(), payload)
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("id does not exist"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
	response := fromDB(author)
	// write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}
