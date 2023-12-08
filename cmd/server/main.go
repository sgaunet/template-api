package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sgaunet/template-api/pkg/config"
	"github.com/sgaunet/template-api/pkg/webserver"
	"github.com/sirupsen/logrus"
)

var version string

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

	// Load config
	cfg, err := config.LoadConfigFromFile(cfgFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	log := initTrace(cfg.DebugLevel)
	w, err := webserver.NewWebServer(cfg, log)
	if err != nil {
		log.Errorf("Error creating webserver: %s\n", err)
		os.Exit(1)
	}
	// handle graceful shutdown
	sigs := make(chan os.Signal, 5)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("Shutting down...")
		w.Shutdown()
	}()
	w.Start()
}

func initTrace(debugLevel string) *logrus.Logger {
	appLog := logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	appLog.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,
		FullTimestamp:    false,
		DisableTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	appLog.SetOutput(os.Stdout)

	switch debugLevel {
	case "debug":
		appLog.SetLevel(logrus.DebugLevel)
	case "warn":
		appLog.SetLevel(logrus.WarnLevel)
	case "error":
		appLog.SetLevel(logrus.ErrorLevel)
	default:
		appLog.SetLevel(logrus.InfoLevel)
	}
	return appLog
}
