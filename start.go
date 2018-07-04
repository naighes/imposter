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

func start(opts *startOpts) (err error) {
	rawConfig := []byte("{}")
	configPath, err := filepath.Abs(opts.configFile)
	if err == nil {
		rawConfig, err = ioutil.ReadFile(opts.configFile)
		if err != nil {
			rawConfig = []byte("{}")
		}
	}
	config, err := ParseConfig(rawConfig)
	if err != nil {
		return fmt.Errorf("could not parse configuration: %v\n", err)
	}
	if opts.configFile != "stdin" {
		log.Printf("read configuration from file '%s'", configPath)
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
	done := make(chan bool)
	log.Printf("starting server listening on port %d", opts.port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %s: %v\n", listenAddr, err)
	}
	<-done
	return nil
}
