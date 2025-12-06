package authors

import (
	"strconv"
	"strings"

	"github.com/sgaunet/template-api/internal/apperror"
)

// Author represents a domain author with business logic.
type Author struct {
	ID   int64
	Name string
	Bio  string
}

// AuthorName value object with validation.
type AuthorName string

// Name length constraints for validation.
const (
	MinNameLength = 5
	MaxNameLength = 20
)

// NewAuthorName creates and validates an author name.
func NewAuthorName(name string) (AuthorName, error) {
	trimmed := strings.TrimSpace(name)

	if len(trimmed) < MinNameLength {
		return "", apperror.NewValidationError(
			"Author name too short",
			map[string]string{
				"field": "name",
				"min":   strconv.Itoa(MinNameLength),
				"value": strconv.Itoa(len(trimmed)),
			},
		)
	}

	if len(trimmed) > MaxNameLength {
		return "", apperror.NewValidationError(
			"Author name too long",
			map[string]string{
				"field": "name",
				"max":   strconv.Itoa(MaxNameLength),
				"value": strconv.Itoa(len(trimmed)),
			},
		)
	}

	return AuthorName(trimmed), nil
}

func (n AuthorName) String() string {
	return string(n)
}

// AuthorBio value object.
type AuthorBio string

// MaxBioLength is the maximum length for author bio.
const MaxBioLength = 500

// NewAuthorBio creates and validates an author bio.
func NewAuthorBio(bio string) (AuthorBio, error) {
	trimmed := strings.TrimSpace(bio)

	if len(trimmed) > MaxBioLength {
		return "", apperror.NewValidationError(
			"Author bio too long",
			map[string]string{
				"field": "bio",
				"max":   strconv.Itoa(MaxBioLength),
				"value": strconv.Itoa(len(trimmed)),
			},
		)
	}

	return AuthorBio(trimmed), nil
}

func (b AuthorBio) String() string {
	return string(b)
}

// NewAuthor creates a new author with validation.
func NewAuthor(name AuthorName, bio AuthorBio) *Author {
	return &Author{
		Name: name.String(),
		Bio:  bio.String(),
	}
}

// CreateAuthorRequest is the request to create an author.
type CreateAuthorRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

// Validate validates the create author request.
func (r *CreateAuthorRequest) Validate() error {
	if _, err := NewAuthorName(r.Name); err != nil {
		return err
	}
	if _, err := NewAuthorBio(r.Bio); err != nil {
		return err
	}
	return nil
}

// ToAuthor converts request to domain author.
func (r *CreateAuthorRequest) ToAuthor() (*Author, error) {
	name, err := NewAuthorName(r.Name)
	if err != nil {
		return nil, err
	}

	bio, err := NewAuthorBio(r.Bio)
	if err != nil {
		return nil, err
	}

	return NewAuthor(name, bio), nil
}

// AuthorResponse is the response format.
type AuthorResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

// ToResponse converts domain author to response.
func (a *Author) ToResponse() *AuthorResponse {
	return &AuthorResponse{
		ID:   a.ID,
		Name: a.Name,
		Bio:  a.Bio,
	}
}
