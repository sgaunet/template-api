package database

import "context"

type StubOK struct {
}

func (s StubOK) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	return Author{}, nil
}

func (s StubOK) DeleteAuthor(ctx context.Context, id int64) error {
	return nil
}

func (s StubOK) GetAuthor(ctx context.Context, id int64) (Author, error) {
	return Author{}, nil
}

func (s StubOK) ListAuthors(ctx context.Context) ([]Author, error) {
	return []Author{}, nil
}

func (s StubOK) PartialUpdateAuthor(ctx context.Context, arg PartialUpdateAuthorParams) (Author, error) {
	return Author{}, nil
}

func (s StubOK) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	return Author{}, nil
}
