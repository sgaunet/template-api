// Package service implements the business logic for the author service.
package service

import (
	"context"
	"fmt"

	"github.com/sgaunet/template-api/internal/repository"
)

//go:generate go tool github.com/matryer/moq -out mock_test.go -pkg service_test ../../../internal/repository Querier

// AuthorService is the Authors Service Layer.
type AuthorService struct {
	queries repository.Querier
}

// NewService creates a new authors service.
func NewService(queries repository.Querier) *AuthorService {
	return &AuthorService{queries: queries}
}

// Create creates a new author.
func (s *AuthorService) Create(ctx context.Context, name string, bio string) (repository.Author, error) {
	author, err := s.queries.CreateAuthor(ctx, repository.CreateAuthorParams{
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		return author, fmt.Errorf("could not create author: %w", err)
	}
	return author, nil
}

// Get returns an author by id.
func (s *AuthorService) Get(ctx context.Context, authorID int64) (*repository.Author, error) {
	// Get author
	author, err := s.queries.GetAuthor(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("could not get author: %w", err)
	}
	return &author, nil
}

// Delete deletes an author by id.
func (s *AuthorService) Delete(ctx context.Context, authorID int64) error {
	if err := s.queries.DeleteAuthor(ctx, authorID); err != nil {
		return fmt.Errorf("could not delete author: %w", err)
	}
	return nil
}

// List returns all authors.
func (s *AuthorService) List(ctx context.Context) ([]repository.Author, error) {
	authors, err := s.queries.ListAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list authors: %w", err)
	}
	return authors, nil
}

// FullUpdate updates an author by id
// func (s *AuthorService) FullUpdate(ctx context.Context, ) (*database.Author, error) {
// 	payload := database.UpdateAuthorParams{}
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		_, _ = w.Write([]byte(err.Error()))
// 		return
// 	}
// 	defer r.Body.Close()
// 	err = json.Unmarshal(body, &payload)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		_, _ = w.Write([]byte(err.Error()))
// 		return
// 	}
// 	author, err := s.queries.UpdateAuthor(context.Background(), payload)
// 	if err != nil {
// 		if err != nil {
// 			if errors.Is(err, sql.ErrNoRows) {
// 				w.WriteHeader(http.StatusBadRequest)
// 				_, _ = w.Write([]byte("id does not exist"))
// 				return
// 			}
// 			w.WriteHeader(http.StatusInternalServerError)
// 			_, _ = w.Write([]byte(err.Error()))
// 			return
// 		}
// 	}
// 	response := fromDB(author)
// 	// write response
// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		_, _ = w.Write([]byte(err.Error()))
// 	}
// }
