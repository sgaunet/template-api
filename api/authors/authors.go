package authors

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sgaunet/template-api/internal/database"
)

type Service struct {
	queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *fiber.App) {
	router.Post("/authors", s.Create)
	router.Get("/authors/:id", s.Get)
	router.Put("/authors/:id", s.FullUpdate)
	router.Delete("/authors/:id", s.Delete)
	router.Get("/authors", s.List)
}

type apiAuthor struct {
	ID   int64
	Name string `json:"name,omitempty" binding:"required,max=32"`
	Bio  string `json:"bio,omitempty" binding:"required"`
}

func (s *Service) Create(c *fiber.Ctx) error {
	payload := struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	// Create author
	params := database.CreateAuthorParams{
		Name: payload.Name,
		Bio:  payload.Bio,
	}
	author, err := s.queries.CreateAuthor(context.Background(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Build response
	response := fromDB(author)
	c.Status(http.StatusCreated)
	return c.JSON(response)
}

func fromDB(author database.Author) *apiAuthor {
	return &apiAuthor{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}

func (s *Service) Get(c *fiber.Ctx) error {
	// Parse request
	id := c.Query("id")
	if id == "" {
		return fiber.NewError(fiber.StatusInternalServerError, errors.New("id is required").Error())
	}
	// convert id to int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Get author
	author, err := s.queries.GetAuthor(context.Background(), idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, errors.New("id does nto exist").Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Build response
	response := fromDB(author)
	return c.JSON(response)
}

func (s *Service) Delete(c *fiber.Ctx) error {
	// Parse request
	id := c.Query("id")
	if id == "" {
		return fiber.NewError(fiber.StatusInternalServerError, errors.New("id is required").Error())
	}
	// convert id to int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Delete author
	if err := s.queries.DeleteAuthor(context.Background(), idInt); err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, errors.New("id does nto exist").Error())
		}
		return fiber.NewError(fiber.StatusServiceUnavailable, errors.New("???").Error())
	}
	return c.SendStatus(http.StatusOK)
}

func (s *Service) List(c *fiber.Ctx) error {
	// List authors
	authors, err := s.queries.ListAuthors(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, errors.New("problem").Error())
	}
	// if len(authors) == 0 {
	// 	return c.SendStatus(http.StatusNotFound)
	// }
	// Build response
	var response []*apiAuthor
	for _, author := range authors {
		response = append(response, fromDB(author))
	}
	return c.JSON(response)
}

func (s *Service) FullUpdate(c *fiber.Ctx) error {
	payload := struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	// Update author
	params := database.UpdateAuthorParams{
		ID:   payload.ID,
		Name: payload.Name,
		Bio:  payload.Bio,
	}
	author, err := s.queries.UpdateAuthor(context.Background(), params)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, errors.New("id does nto exist").Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	// Build response
	response := fromDB(author)
	return c.JSON(response)
}
