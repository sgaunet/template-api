package authors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sgaunet/template-api/api/authors"
	"github.com/sgaunet/template-api/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOK(t *testing.T) {
	mock := &QuerierMock{
		DeleteAuthorFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	svc := authors.NewService(mock)
	Router := chi.NewRouter()
	Router.Delete("/author/{id}", svc.Delete)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/author/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteWithWrongID(t *testing.T) {
	mock := &QuerierMock{
		DeleteAuthorFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	svc := authors.NewService(mock)
	Router := chi.NewRouter()
	Router.Delete("/author/{id}", svc.Delete)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/author/WrongID", nil)
	if err != nil {
		t.Fatal(err)
	}
	Router.ServeHTTP(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestDeleteWithError(t *testing.T) {
	mock := &QuerierMock{
		DeleteAuthorFunc: func(ctx context.Context, id int64) error {
			return errors.New("error")
		},
	}
	svc := authors.NewService(mock)
	Router := chi.NewRouter()
	Router.Delete("/author/{id}", svc.Delete)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/author/WrongID", nil)
	if err != nil {
		t.Fatal(err)
	}
	Router.ServeHTTP(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestList(t *testing.T) {
	mock := &QuerierMock{
		ListAuthorsFunc: func(ctx context.Context) ([]database.Author, error) {
			return []database.Author{
				{
					ID:   1,
					Name: "John Doe",
					Bio:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				},
			}, nil
		},
	}
	// Create a new service
	svc := authors.NewService(mock)
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc.List(w, r)
		}))
	defer ts.Close()
	// get request
	client := ts.Client()
	resp, err := client.Get(ts.URL + "/authors") //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var authors []database.Author
	err = json.NewDecoder(resp.Body).Decode(&authors)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equalf(t, 1, len(authors), "number of authors")
	// Verify, if the status code is as expected
	assert.Equalf(t, http.StatusOK, resp.StatusCode, "status code")
}

func TestCreate(t *testing.T) {
	mock := &QuerierMock{
		CreateAuthorFunc: func(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) {
			return database.Author{
				ID:   1,
				Name: arg.Name,
				Bio:  arg.Bio,
			}, nil
		},
	}
	// Create a new service
	svc := authors.NewService(mock)
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc.Create(w, r)
		}))
	defer ts.Close()

	// get request
	client := ts.Client()

	payload := struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}{
		Name: "John Doe",
		Bio:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Post(ts.URL+"/authors", "application/json", bytes.NewReader(data)) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	// unmarshal response
	var author database.Author
	err = json.NewDecoder(resp.Body).Decode(&author)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equalf(t, "John Doe", author.Name, "author name")
	assert.Equalf(t, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", author.Bio, "author bio")
	assert.NotEqualf(t, int64(0), author.ID, "author id")
	// Verify, if the status code is as expected
	assert.Equalf(t, http.StatusCreated, resp.StatusCode, "status code")
}

func TestCreateNameTooShort(t *testing.T) {
	mock := &QuerierMock{
		CreateAuthorFunc: func(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) {
			return database.Author{
				ID:   1,
				Name: arg.Name,
				Bio:  arg.Bio,
			}, nil
		},
	}
	// Create a new service
	svc := authors.NewService(mock)
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc.Create(w, r)
		}))
	defer ts.Close()

	// get request
	client := ts.Client()

	payload := struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}{
		Name: "Joh",
		Bio:  "Lor.",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Post(ts.URL+"/authors", "application/json", bytes.NewReader(data)) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	// Verify, if the status code is as expected
	assert.Equalf(t, http.StatusBadRequest, resp.StatusCode, "status code")
}

func TestGetOK(t *testing.T) {
	mock := &QuerierMock{
		GetAuthorFunc: func(ctx context.Context, id int64) (database.Author, error) {
			return database.Author{
				ID:   1,
				Name: "John Doe",
				Bio:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			}, nil
		},
	}
	svc := authors.NewService(mock)
	Router := chi.NewRouter()
	Router.Get("/author/{id}", svc.Get)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/author/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}
