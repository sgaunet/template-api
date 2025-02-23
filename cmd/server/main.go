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

//go:generate sqlc generate

const channelSignalSize = 5

var version string = "development"

func printVersion() {
	fmt.Printf("%s\n", version)
}

func main() {
	var (
		err         error
		cfg         config.Config
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
	if cfgFile == "" {
		cfg = config.LoadConfigFromEnvVar()
	} else {
		cfg, err = config.LoadConfigFromFile(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading YAML file: %s\n", err)
			os.Exit(1)
		}
	}

	if !cfg.IsValid() {
		fmt.Fprintf(os.Stderr, "Invalid configuration\n")
		fmt.Fprintf(os.Stderr, "DBDSN: %s\n", cfg.DBDSN)
		fmt.Fprintf(os.Stderr, "REDISDSN: %s\n", cfg.RedisDSN)
		os.Exit(1)
	}

	// wait for database
	waitCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	err = database.WaitForDB(waitCtx, cfg.DBDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error waiting for database: %s\n", err)
		os.Exit(1)
	}
	cancel()

	// init database connection
	pg, err := database.NewPostgres(cfg.DBDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to database failed: %s\n", err.Error())
		os.Exit(1)
	}
	// init database (create tables, etc...)
	err = pg.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %s\n", err)
		pg.Close()
		os.Exit(1)
	}
	// init redis connection
	// _, err = initRedisConnection(cfg.RedisDSN)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error initializing redis: %s\n", err)
	// 	pg.Close()
	// 	os.Exit(1)
	// }

	// init services
	authorsQueries := repository.New(pg.GetDB())
	authorSvc := service.NewService(authorsQueries)

	// init webserver
	w, err := webserver.NewWebServer(authorSvc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating webserver: %s\n", err)
		os.Exit(1)
	}

	// handle graceful shutdown
	sigs := make(chan os.Signal, channelSignalSize)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = w.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error starting webserver: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	<-sigs
	err = w.Shutdown(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error shutting down webserver: %s\n", err.Error())
		os.Exit(1)
	}
}

//nolint:unused
func initRedisConnection(redisdsn string) (*redis.Client, error) {
	var err error
	d, err := dsn.New(redisdsn)
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%s:%s", d.GetHost(), d.GetPort("6379"))
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	return redisClient, err
}
