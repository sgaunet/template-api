package database_test

import (
	"log"

	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresqlC testcontainers.Container
var pgdsn string

func setup() {
	// Do something here.
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
	postgresqlC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	endpoint, err := postgresqlC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err.Error())
	}
	pgdsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", "postgres", "password", endpoint, "postgres")
}

func teardown() {
	// Do something here.
	defer func() {
		if err := postgresqlC.Terminate(context.Background()); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}
