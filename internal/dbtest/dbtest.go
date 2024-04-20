package dbtest

import (
	"log"

	"context"
	"fmt"

	"github.com/sgaunet/template-api/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB is a struct that contains the test database container and the DSN to connect to it.
type TestDB struct {
	postgresqlC testcontainers.Container
	pgDSN       string
}

// NewTestDB creates a new Test database (need docker to be installed and running).
// If waitDbForBeingReady is true, it will wait for the database to be ready to accept connections.
// It returns a pointer to a TestDB struct.
// The TestDB struct contains the test database container and the DSN to connect to it.
// Don't forget to call the Teardown method to clean up the resources.
func NewTestDB(waitDbForBeingReady bool) *TestDB {
	newpgDB := &TestDB{}
	ctx := context.Background()
	var err error
	// req := testcontainers.ContainerRequest{
	//     Image:        "redis:latest",
	//     ExposedPorts: []string{"6379/tcp"},
	//     WaitingFor:   wait.ForLog("Ready to accept connections"),
	// }
	// redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	//     ContainerRequest: req,
	//     Started:          true,
	// })
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16.2",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
		},
	}
	newpgDB.postgresqlC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	endpoint, err := newpgDB.postgresqlC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err.Error())
	}
	newpgDB.pgDSN = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", "postgres", "password", endpoint, "postgres")
	if waitDbForBeingReady {
		err = database.WaitForDB(ctx, newpgDB.pgDSN)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return newpgDB
}

// Teardown cleans up the resources.
func (t *TestDB) Teardown() {
	defer func() {
		if err := t.postgresqlC.Terminate(context.Background()); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

// GetDSN returns the DSN to connect to the test database.
func (t *TestDB) GetDSN() string {
	return t.pgDSN
}
