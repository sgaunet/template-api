package service_test

import (
	"context"
	"testing"

	"github.com/sgaunet/template-api/pkg/authors/service"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOK(t *testing.T) {
	var argID int64
	mock := &QuerierMock{
		DeleteAuthorFunc: func(ctx context.Context, id int64) error {
			argID = id
			return nil
		},
	}
	svc := service.NewService(mock)
	svc.Delete(context.Background(), 1)
	assert.Equal(t, int64(1), argID)
}
