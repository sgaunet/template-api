// Package domain contains the domain models for the application.
package domain

import (
	"github.com/sgaunet/template-api/internal/repository"
)

// AuthorRequestBody is the request body for the Author entity.
type AuthorRequestBody struct {
	ID   int64  `json:"-"`
	Name string `json:"name,omitempty" validate:"min=5,max=20"`
	Bio  string `json:"bio,omitempty"`
}

// AuthorResponseBody is the response body for the Author entity.
type AuthorResponseBody struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty" validate:"min=5,max=20"`
	Bio  string `json:"bio,omitempty"`
}

// FromDB converts a repository.Author to an AuthorResponseBody.
func FromDB(author repository.Author) *AuthorResponseBody {
	return &AuthorResponseBody{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}
