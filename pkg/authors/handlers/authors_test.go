package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sgaunet/template-api/pkg/authors/handlers"
	"github.com/sgaunet/template-api/pkg/authors/service"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOK(t *testing.T) {
	mock := &QuerierMock{
		DeleteAuthorFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	svc := service.NewService(mock)
	h := handlers.NewAuthorsHandlers(svc)
	Router := chi.NewRouter()
	Router.Delete("/author/{uuid}", h.Delete)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/author/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	// assert that mock has been called
	assert.Equal(t, int64(1), mock.calls.DeleteAuthor[0].ID)
}
