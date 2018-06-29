package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Args struct {
	port       int
	configFile string
}

func ParseArgs(args []string) (*Args, error) {
	r := Args{port: 8080}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--port":
			i = i + 1
			if len(args)-1 < i {
				return nil, fmt.Errorf("expected value for token '%s'", arg)
			}
			v, err := strconv.Atoi(args[i])
			(&r).port = v
			if err != nil {
				return nil, fmt.Errorf("unexpected value for token '%s': expected 'int'", arg)
			}
		case "--config-file":
			i = i + 1
			if len(args)-1 < i {
				return nil, fmt.Errorf("expected value for token '%s'", arg)
			}
			(&r).configFile = args[i]
		default:
			return nil, fmt.Errorf("flag '%s'is not recognized", arg)
		}
	}
	return &r, nil
}

func main() {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config, err := ioutil.ReadFile(args.configFile)
	if err != nil {
		config = []byte("[]")
	}
	defs, err := ParseMatchDef(config)
	if err != nil {
		fmt.Printf("could not parse configuration: %v\n", err)
		os.Exit(1)
	}
	router, err := NewRegexHandler(defs)
	if err != nil {
		fmt.Printf("could not load configuration: %v\n", err)
		os.Exit(1)
	}
	listenAddr := fmt.Sprintf("localhost:%d", args.port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
	done := make(chan bool)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("could not listen on %s: %v\n", listenAddr, err)
	}
	<-done
}
