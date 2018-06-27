package main

import (
	"fmt"
	"net/http"
	"os"
)

// TODO: lower and upper case check for function names

// FEATURES
// [BIG] generate from swagger
// [BIG] record all requests

func main() {
	j := []byte(`[{"pattern": "^/myfile$", "response": {"body": "${file(./main.go)}", "headers": {"Content-Type": "text/plain; charset=utf-8"}, "status_code": 200}}, {"pattern": "^/[a-z]+$", "response": {"body": "${text(Hello, string!)}", "headers": {"Content-Type": "text/plain; charset=utf-8"}, "status_code": 200}}, {"pattern": "^/[0-9]+$", "response": {"body": "Hello, number!", "headers": {"Content-Type": "text/plain; charset=utf-8"}, "status_code": 404}}]`)
	defs, err := ParseMatchDef(j)
	if err != nil {
		// ERROR: exit program
		fmt.Println(err)
		os.Exit(1)
	}
	router, err := NewRegexHandler(defs)
	if err != nil {
		// ERROR: exit program
		fmt.Println(err)
		os.Exit(1)
	}
	listenAddr := "localhost:5000"
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
	done := make(chan bool)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Could not listen on %s: %v\n", listenAddr, err)
	}
	<-done
}
