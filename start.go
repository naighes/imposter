package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

const defaultPort int = 8080

func startCmd() command {
	fs := flag.NewFlagSet("imposter start", flag.ExitOnError)
	opts := startOpts{port: defaultPort}
	fs.IntVar(&opts.port, "port", defaultPort, "The listening TCP port")
	fs.StringVar(&opts.configFile, "config-file", "stdin", "The configuration file")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return start(&opts)
	}}
}

type startOpts struct {
	port       int
	configFile string
}

func readConfig(configFile string) *Config {
	if configPath, err := filepath.Abs(configFile); err == nil {
		if rawConfig, err := ioutil.ReadFile(configPath); err == nil {
			if config, err := ParseConfig(rawConfig); err == nil {
				log.Printf("read configuration from file '%s'", configPath)
				return config
			}
		}
	}
	return &Config{}
}

func start(opts *startOpts) (err error) {
	config := readConfig(opts.configFile)
	router, err := NewRegexHandler(config)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v\n", err)
	}
	listenAddr := fmt.Sprintf("localhost:%d", opts.port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
	done := make(chan bool)
	log.Printf("starting server listening on port %d", opts.port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %s: %v\n", listenAddr, err)
	}
	<-done
	return nil
}
