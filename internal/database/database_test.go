package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sgaunet/template-api/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// TestWaitDBHandlesCtx tests the WaitForDB function with a context
// that is cancelled before the database is ready
// This test must be the first one to run
func TestWaitDBHandlesCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	err := database.WaitForDB(ctx, pgdsn)
	assert.NotNil(t, err)
}
func TestWaitDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := database.WaitForDB(ctx, pgdsn)
	assert.Nil(t, err)
}

// TestInitDB tests the InitDB function
func TestInitDB(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)

	err = pg.InitDB()
	assert.Nil(t, err)
}

func TestGetDB(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	db := pg.GetDB()
	assert.NotNil(t, db)
}

func TestCreateAuthor(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := database.New(pg.DB)
	var author database.Author
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)
}

func TestGetAuthor(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := database.New(pg.DB)
	var author database.Author
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// Get author
	author, err = q.GetAuthor(context.Background(), author.ID)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)
	assert.Equal(t, "John Doe", author.Name)
}

func TestDeleteAuthor(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := database.New(pg.DB)
	var author database.Author
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// Delete author
	err = q.DeleteAuthor(context.Background(), author.ID)
	assert.Nil(t, err)

	// Try to Get author
	author, err = q.GetAuthor(context.Background(), author.ID)
	assert.NotNil(t, err)
}

func TestUpdateAuthor(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := database.New(pg.DB)
	var author, updatedAuthor database.Author
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple test",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// Delete author
	updatedAuthor, err = q.UpdateAuthor(context.Background(), database.UpdateAuthorParams{
		ID:   author.ID,
		Name: "John Doe Jr",
		Bio:  "A simple ...............",
	})
	assert.Nil(t, err)
	assert.Equal(t, "John Doe Jr", updatedAuthor.Name)
	assert.Equal(t, "A simple ...............", updatedAuthor.Bio)
}

func TestListAuthor(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := database.New(pg.DB)
	var author database.Author
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe 2",
		Bio:  "A simple 2",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// List authors
	authors, err := q.ListAuthors(context.Background())
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(authors))
}

func TestUpdatePartialAuthor(t *testing.T) {
	database.WaitForDB(context.Background(), pgdsn)
	pg, err := database.NewPostgres(pgdsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := database.New(pg.DB)
	var author database.Author
	author, err = q.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// Update Partial
	author, err = q.PartialUpdateAuthor(context.Background(), database.PartialUpdateAuthorParams{
		ID:         author.ID,
		UpdateBio:  true,
		Bio:        "A simple test",
		UpdateName: true,
		Name:       "John Doe Jr",
	})
	assert.Nil(t, err)
	assert.Equal(t, "John Doe Jr", author.Name)
	assert.Equal(t, "A simple test", author.Bio)
}
