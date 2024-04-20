package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/sgaunet/template-api/internal/database"
	"github.com/sgaunet/template-api/internal/dbtest"
	"github.com/sgaunet/template-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

var testdb *dbtest.TestDB

func TestMain(m *testing.M) {
	testdb = dbtest.NewTestDB(true)
	code := m.Run()
	testdb.Teardown()
	os.Exit(code)
}

func TestCreateAuthor(t *testing.T) {
	err := database.WaitForDB(context.Background(), testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := repository.New(pg.DB)
	var author repository.Author
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)
}

func TestGetAuthor(t *testing.T) {
	err := database.WaitForDB(context.Background(), testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := repository.New(pg.DB)
	var author repository.Author
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
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
	err := database.WaitForDB(context.Background(), testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := repository.New(pg.DB)
	var author repository.Author
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
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
	if err := database.WaitForDB(context.Background(), testdb.GetDSN()); err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := repository.New(pg.DB)
	var author, updatedAuthor repository.Author
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple test",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// Delete author
	updatedAuthor, err = q.UpdateAuthor(context.Background(), repository.UpdateAuthorParams{
		ID:   author.ID,
		Name: "John Doe Jr",
		Bio:  "A simple ...............",
	})
	assert.Nil(t, err)
	assert.Equal(t, "John Doe Jr", updatedAuthor.Name)
	assert.Equal(t, "A simple ...............", updatedAuthor.Bio)
}

func TestListAuthor(t *testing.T) {
	if err := database.WaitForDB(context.Background(), testdb.GetDSN()); err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := repository.New(pg.DB)
	var author repository.Author
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
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
	if err := database.WaitForDB(context.Background(), testdb.GetDSN()); err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	assert.NotNil(t, pg.DB)
	assert.Nil(t, err)
	err = pg.InitDB()
	assert.Nil(t, err)

	// Create author
	q := repository.New(pg.DB)
	var author repository.Author
	author, err = q.CreateAuthor(context.Background(), repository.CreateAuthorParams{
		Name: "John Doe",
		Bio:  "A simple",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, author.ID)

	// Update Partial
	author, err = q.PartialUpdateAuthor(context.Background(), repository.PartialUpdateAuthorParams{
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
