package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/sgaunet/dsn"
	"github.com/sgaunet/template-api/api/authors"
	"github.com/sgaunet/template-api/internal/database"
	"github.com/sgaunet/template-api/pkg/config"
	"github.com/sgaunet/template-api/pkg/webserver"
)

const channelSignalSize = 5

var version string = "development"

func printVersion() {
	fmt.Printf("%s\n", version)
}

func main() {
	var (
		err         error
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

	cfg, err := config.LoadConfigFromFileOrEnvVar(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML file: %s\n", err)
		os.Exit(1)
	}
	if !cfg.IsValid() {
		fmt.Fprintf(os.Stderr, "Invalid configuration\n")
		os.Exit(1)
	}

	// init database connection
	pg, err := database.NewPostgres(cfg.DBDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration\n")
		os.Exit(1)
	}

	// init database
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

	queries := database.New(pg.GetDB())
	authorSvc := authors.NewService(queries)
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
