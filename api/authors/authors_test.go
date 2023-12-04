package authors_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sgaunet/template-api/api/authors"
	"github.com/sgaunet/template-api/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	stub := database.StubOK{}
	// Create a new service
	svc := authors.NewService(stub)
	app := fiber.New()
	app.Get("/authors", svc.List)
	// Call the service
	req := httptest.NewRequest(http.MethodGet, "/authors", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	// Verify, if the status code is as expected
	assert.Equalf(t, fiber.StatusOK, resp.StatusCode, "status code")
}

func TestCreate(t *testing.T) {
	stub := database.StubOK{}
	// Create a new service
	svc := authors.NewService(stub)
	app := fiber.New()
	app.Post("/authors", svc.Create)
	payload := struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}{
		Name: "John Doe",
		Bio:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(data)
	req := httptest.NewRequest(http.MethodPost, "/authors", reader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	// Verify, if the status code is as expected
	assert.Equalf(t, fiber.StatusCreated, resp.StatusCode, "status code")
}
