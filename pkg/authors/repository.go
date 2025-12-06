package authors

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sgaunet/template-api/internal/apperror"
	"github.com/sgaunet/template-api/internal/repository"
)

// Repository defines the interface for author data access.
type Repository interface {
	Create(ctx context.Context, author *Author) (*Author, error)
	GetByID(ctx context.Context, id int64) (*Author, error)
	List(ctx context.Context) ([]*Author, error)
	Delete(ctx context.Context, id int64) error
}

// repositoryImpl wraps sqlc-generated queries.
type repositoryImpl struct {
	queries repository.Querier
}

// NewRepository creates a new author repository.
func NewRepository(queries repository.Querier) Repository {
	return &repositoryImpl{queries: queries}
}

func (r *repositoryImpl) Create(ctx context.Context, author *Author) (*Author, error) {
	dbAuthor, err := r.queries.CreateAuthor(ctx, repository.CreateAuthorParams{
		Name: author.Name,
		Bio:  author.Bio,
	})
	if err != nil {
		return nil, apperror.NewInternalError(err)
	}

	return &Author{
		ID:   dbAuthor.ID,
		Name: dbAuthor.Name,
		Bio:  dbAuthor.Bio,
	}, nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id int64) (*Author, error) {
	dbAuthor, err := r.queries.GetAuthor(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Author not found")
		}
		return nil, apperror.NewInternalError(err)
	}

	return &Author{
		ID:   dbAuthor.ID,
		Name: dbAuthor.Name,
		Bio:  dbAuthor.Bio,
	}, nil
}

func (r *repositoryImpl) List(ctx context.Context) ([]*Author, error) {
	dbAuthors, err := r.queries.ListAuthors(ctx)
	if err != nil {
		return nil, apperror.NewInternalError(err)
	}

	authors := make([]*Author, len(dbAuthors))
	for i, dbAuthor := range dbAuthors {
		authors[i] = &Author{
			ID:   dbAuthor.ID,
			Name: dbAuthor.Name,
			Bio:  dbAuthor.Bio,
		}
	}

	return authors, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id int64) error {
	if err := r.queries.DeleteAuthor(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Author not found")
		}
		return apperror.NewInternalError(err)
	}
	return nil
}
