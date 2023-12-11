package webserver_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"context"
	"fmt"

	"github.com/sgaunet/template-api/pkg/config"
	"github.com/sgaunet/template-api/pkg/webserver"
	"github.com/stretchr/testify/assert"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresqlC testcontainers.Container
var pgdsn string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

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
		Image:        "postgres:15.4",
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

func TestGetAuthors(t *testing.T) {
	time.Sleep(time.Second * 1)
	w, err := webserver.NewWebServer(&config.Config{
		DbDSN: pgdsn,
		// DbDSN: "postgres://postgres:password@localhost:15432/postgres?sslmode=disable",
	})
	if err != nil {
		t.Fatal(err)
	}
	go w.Start()
	defer w.Shutdown()

	req, err := http.NewRequest("GET", "http://localhost:3000/authors", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	// Verify, if the status code is as expected
	assert.Equalf(t, http.StatusOK, resp.StatusCode, "Status code %v", resp.StatusCode)
}

func TestGetAuthorsWithValues(t *testing.T) {
	time.Sleep(time.Second * 1)
	w, err := webserver.NewWebServer(&config.Config{
		DbDSN: pgdsn,
		// DbDSN: "postgres://postgres:password@localhost:15432/postgres?sslmode=disable",
	})
	if err != nil {
		t.Fatal(err)
	}
	go w.Start()
	defer w.Shutdown()

	// add a new author
	data := []byte(`{"name":"John Doe","bio":"A great author"}`)
	bodyReader := bytes.NewReader(data)
	req, err := http.NewRequest("POST", "http://localhost:3000/authors", bodyReader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, resp.StatusCode)
	}

	// get all authors
	req, err = http.NewRequest("GET", "http://localhost:3000/authors", nil)
	if err != nil {
		t.Error(err)
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp2.Body.Close()

	// json decode the response body
	var authors []map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&authors)
	if err != nil {
		t.Error(err)
	}
	// Verify, if the status code is as expected
	assert.Equalf(t, http.StatusOK, resp2.StatusCode, "Status code %v", resp2.StatusCode)
	// Verify, if the response body is as expected
	// assert.Equalf(t, 1, len(authors), "Number of authors %v", len(authors))

	// Verify, if the status code is as expected
	assert.Equalf(t, http.StatusCreated, resp.StatusCode, "Status code %v", resp.StatusCode)
}
