// Package authors provides the authors domain logic.
package authors

import (
	"context"
	"fmt"

	"github.com/sgaunet/template-api/internal/apperror"
)

// Service provides author business logic.
type Service interface {
	Create(ctx context.Context, req *CreateAuthorRequest) (*Author, error)
	GetByID(ctx context.Context, id int64) (*Author, error)
	List(ctx context.Context) ([]*Author, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

// NewService creates a new author service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateAuthorRequest) (*Author, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Convert to domain model
	author, err := req.ToAuthor()
	if err != nil {
		return nil, err
	}

	// Business logic can go here (e.g., check for duplicates, apply business rules)

	// Persist
	return s.repo.Create(ctx, author)
}

func (s *service) GetByID(ctx context.Context, id int64) (*Author, error) {
	if id <= 0 {
		return nil, apperror.NewValidationError(
			"Invalid author ID",
			map[string]string{"field": "id", "value": fmt.Sprintf("%d", id)},
		)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]*Author, error) {
	return s.repo.List(ctx)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperror.NewValidationError(
			"Invalid author ID",
			map[string]string{"field": "id", "value": fmt.Sprintf("%d", id)},
		)
	}

	return s.repo.Delete(ctx, id)
}
