package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const defaultPort int = 8080

func startCmd() command {
	fs := flag.NewFlagSet("imposter start", flag.ExitOnError)
	opts := startOpts{port: defaultPort}
	fs.IntVar(&opts.port, "port", defaultPort, "The listening TCP port")
	fs.StringVar(&opts.configFile, "config-file", "stdin", "The configuration file")
	flag.DurationVar(&opts.wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return start(&opts)
	}}
}

type startOpts struct {
	port       int
	configFile string
	wait       time.Duration
}

func readConfig(configFile string) (*Config, error) {
	var err error
	var configPath string
	var rawConfig []byte
	var config *Config
	if configPath, err = filepath.Abs(configFile); err == nil {
		if rawConfig, err = ioutil.ReadFile(configPath); err == nil {
			if config, err = ParseConfig(rawConfig); err == nil {
				log.Printf("read configuration from file '%s'\n", configPath)
				return config, nil
			}
		}
	}
	return &Config{}, err
}

func start(opts *startOpts) error {
	config, err := readConfig(opts.configFile)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v\n", err)
	}
	router, err := NewRegexHandler(config)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v\n", err)
	}
	listenAddr := fmt.Sprintf("localhost:%d", opts.port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
	c := make(chan os.Signal, 1)
	log.Printf("starting imposter instance listening on port %d...\n", opts.port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("could not listen on %s: %v\n", listenAddr, err)
		}
	}()
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), opts.wait)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("imposter is shutting down...")
	return nil
}
