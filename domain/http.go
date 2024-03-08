package domain

import (
	"github.com/sgaunet/template-api/pkg/authors/repository"
)

// apiAuthor is the author representation for the API
type AuthorRequestBody struct {
	ID   int64
	Name string `json:"name,omitempty" validate:"min=5,max=20"`
	Bio  string `json:"bio,omitempty"`
}

// apiAuthor is the author representation for the API
type AuthorResponseBody struct {
	ID   int64
	Name string `json:"name,omitempty" validate:"min=5,max=20"`
	Bio  string `json:"bio,omitempty"`
}

func FromDB(author repository.Author) *AuthorResponseBody {
	return &AuthorResponseBody{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}
