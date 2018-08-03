package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

const defaultPort int = 8080

func startCmd() command {
	fs := flag.NewFlagSet("imposter start", flag.ExitOnError)
	opts := startOpts{port: defaultPort}
	fs.IntVar(&opts.port, "port", defaultPort, "The listening TCP port")
	fs.StringVar(&opts.configFile, "config-file", "stdin", "The configuration file")
	fs.StringVar(&opts.rawTLSCertFileList, "tls-cert-file-list", "", "A comma separated list of X.509 certificates to secure communication")
	fs.StringVar(&opts.rawTLSKeyFileList, "tls-key-file-list", "", "A comma separated list of private key files corresponding to the X.509 certificates")
	fs.DurationVar(&opts.wait, "graceful-timeout", time.Second*15, "The duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	fs.BoolVar(&opts.record, "record", false, "Enable the recording of PUT requests")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return startExec(&opts)
	}}
}

type startOpts struct {
	port               int
	configFile         string
	wait               time.Duration
	rawTLSCertFileList string
	rawTLSKeyFileList  string
	record             bool
}

func (o *startOpts) buildListenAndServe(server *http.Server) (func() error, error) {
	if o.rawTLSCertFileList != "" && o.rawTLSKeyFileList != "" {
		certs := strings.Split(o.rawTLSCertFileList, ",")
		keys := strings.Split(o.rawTLSKeyFileList, ",")
		if len(certs) != len(keys) {
			return nil, fmt.Errorf("the number of X.509 certificates does not match the number of keys")
		}
		if len(certs) == 1 {
			return func() error {
				return server.ListenAndServeTLS(o.rawTLSCertFileList, o.rawTLSKeyFileList)
			}, nil
		}
		cfg := &tls.Config{}
		for index, cert := range certs {
			pair, err := tls.LoadX509KeyPair(cert, keys[index])
			if err != nil {
				return nil, fmt.Errorf("could not load X.509 pair from %s/%s: %v", cert, keys[index], err)
			}
			cfg.Certificates = append(cfg.Certificates, pair)
		}
		cfg.BuildNameToCertificate()
		server.TLSConfig = cfg
		return func() error {
			return server.ListenAndServeTLS("", "")
		}, nil
	}
	return server.ListenAndServe, nil
}

func startExec(opts *startOpts) error {
	config, err := readConfig(opts.configFile)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v", err)
	}
	var store Store
	if opts.record {
		store = newInMemoryStore()
	} else {
		store = nil
	}
	router, err := NewRouter(config, store)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v", err)
	}
	listenAddr := fmt.Sprintf("localhost:%d", opts.port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
	c := make(chan os.Signal, 1)
	listenAndServe, err := opts.buildListenAndServe(server)
	if err != nil {
		return err
	}
	log.Printf("starting imposter instance listening on port %d...\n", opts.port)
	go func() {
		if err := listenAndServe(); err != nil {
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
