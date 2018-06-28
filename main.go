package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// FEATURES
// [BIG] generate from swagger
// [BIG] record all requests

type Args struct {
	port int
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
			p, err := strconv.Atoi(args[i])
			(&r).port = p
			if err != nil {
				return nil, fmt.Errorf("unexpected value for token '%s': expected 'int'", arg)
			}
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
	j := []byte(`[{"pattern": "^/myfile$", "response": {"body": "${file(./main.go)}", "headers": {"Content-Type": "text/plain; charset=utf-8"}, "status_code": 200}}, {"pattern": "^/[a-z]+$", "response": {"body": "${text(Hello, string!)}", "headers": {"Content-Type": "text/plain; charset=utf-8"}, "status_code": 200}}, {"pattern": "^/[0-9]+$", "response": {"body": "Hello, number!", "headers": {"Content-Type": "text/plain; charset=utf-8"}, "status_code": 404}}]`)
	defs, err := ParseMatchDef(j)
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
