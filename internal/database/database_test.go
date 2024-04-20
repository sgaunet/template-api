package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/sgaunet/template-api/internal/database"
	"github.com/sgaunet/template-api/internal/dbtest"
	"github.com/stretchr/testify/assert"
)

var testdb *dbtest.TestDB

// TestWaitDBHandlesCtx tests the WaitForDB function with a context
// that is cancelled before the database is ready
// This test must be the first one to run
func TestWaitDBHandlesCtx(t *testing.T) {
	testdb = dbtest.NewTestDB(false)
	defer testdb.Teardown()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	err := database.WaitForDB(ctx, testdb.GetDSN())
	assert.NotNil(t, err)
}

func TestWaitDB(t *testing.T) {
	testdb = dbtest.NewTestDB(false)
	defer testdb.Teardown()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := database.WaitForDB(ctx, testdb.GetDSN())
	assert.Nil(t, err)
}

// TestInitDB tests the InitDB function
func TestInitDB(t *testing.T) {
	testdb = dbtest.NewTestDB(true)
	defer testdb.Teardown()
	pg, err := database.NewPostgres(testdb.GetDSN())
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
	testdb = dbtest.NewTestDB(true)
	defer testdb.Teardown()
	ctx := context.Background()
	err := database.WaitForDB(ctx, testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	pg, err := database.NewPostgres(testdb.GetDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close()
	db := pg.GetDB()
	assert.NotNil(t, db)
}
