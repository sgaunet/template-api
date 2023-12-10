package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sgaunet/template-api/pkg/config"
	"github.com/sgaunet/template-api/pkg/webserver"
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

	cfg, err := config.LoadConfigFromFileOrEnvVar(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML file: %s\n", err)
		os.Exit(1)
	}
	if !cfg.IsValid() {
		fmt.Fprintf(os.Stderr, "Invalid configuration\n")
		os.Exit(1)
	}

	w, err := webserver.NewWebServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating webserver: %s\n", err)
		os.Exit(1)
	}
	// handle graceful shutdown
	sigs := make(chan os.Signal, 5)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		w.Shutdown()
	}()
	w.Start()
}
