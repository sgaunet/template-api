package dbtest

import (
	"log"

	"context"
	"fmt"

	"github.com/sgaunet/template-api/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	postgresqlC testcontainers.Container
	PgDSN       string
}

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
	newpgDB.PgDSN = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", "postgres", "password", endpoint, "postgres")
	if waitDbForBeingReady {
		err = database.WaitForDB(ctx, newpgDB.PgDSN)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return newpgDB
}

func (t *TestDB) Teardown() {
	// Do something here.
	defer func() {
		if err := t.postgresqlC.Terminate(context.Background()); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}
