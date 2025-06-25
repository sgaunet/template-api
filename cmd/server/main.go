// Command server is the entry point for the API server.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"github.com/sgaunet/template-api/internal/database"
	"github.com/sgaunet/template-api/internal/repository"
	"github.com/sgaunet/template-api/pkg/authors/service"
	"github.com/sgaunet/template-api/pkg/config"
	"github.com/sgaunet/template-api/pkg/webserver"
)

//go:generate go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate -f ../../sqlc.yaml

const (
	channelSignalSize = 5
	waitForDB         = 30 * time.Second
)

var version = "development"

func printVersion() {
	fmt.Printf("%s\n", version)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		cfgFile     string
		versionFlag bool
	)
	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()

	if versionFlag {
		printVersion()
		os.Exit(0)
	}

	// load configuration
	cfg, err := loadConfiguration(cfgFile)
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// init database
	pg, err := initDB(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := pg.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing database: %v\n", err)
		}
	}()

	// init services
	authorsQueries := repository.New(pg.GetDB())
	authorSvc := service.NewService(authorsQueries)

	// init webserver
	w, err := webserver.NewWebServer(authorSvc)
	if err != nil {
		return fmt.Errorf("error creating webserver: %w", err)
	}

	// handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- w.Start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("error starting webserver: %w", err)
		}
	case <-sigs:
	}

	if err := w.Shutdown(context.TODO()); err != nil {
		return fmt.Errorf("error shutting down webserver: %w", err)
	}

	return nil
}

//nolint:unused
func initRedisConnection(redisdsn string) (*redis.Client, error) {
	var err error
	d, err := dsn.New(redisdsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse redis dsn: %w", err)
	}
	addr := fmt.Sprintf("%s:%s", d.GetHost(), d.GetPort("6379"))
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		return nil, fmt.Errorf("could not ping redis: %w", err)
	}
	return redisClient, nil
}

func initDB(cfg config.Config) (*database.Postgres, error) {
	// wait for database
	waitCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(waitForDB))
	defer cancel()
	if err := database.WaitForDB(waitCtx, cfg.DBDSN); err != nil {
		return nil, fmt.Errorf("error waiting for database: %w", err)
	}
	// init database connection
	pg, err := database.NewPostgres(cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}
	// init database (create tables, etc...)
	if err = pg.InitDB(); err != nil {
		if errClose := pg.Close(); errClose != nil {
			fmt.Fprintf(os.Stderr, "Error closing database during initDB failure: %s\n", errClose)
		}
		return nil, fmt.Errorf("error initializing database: %w", err)
	}
	return pg, nil
}

func loadConfiguration(cfgFile string) (config.Config, error) {
	var (
		err error
		cfg config.Config
	)
	cfg, err = config.Load(cfgFile)
	if err != nil {
		return config.Config{}, fmt.Errorf("error loading configuration: %w", err)
	}

		if err := cfg.Validate(); err != nil {
		return config.Config{}, fmt.Errorf("invalid configuration: %w", err)
	}
	return cfg, nil
}
