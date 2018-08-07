package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Logger defines a general abstraction for the logging of HTTP requests.
type Logger interface {
	log(r *http.Request)
}

// DefaultLogger is a basic implementation of Logger.
type DefaultLogger struct {
}

func (l *DefaultLogger) log(r *http.Request) {
	if r == nil {
		return
	}
	var b bytes.Buffer
	for k := range r.Header {
		b.WriteString(fmt.Sprintf("%s: %s\n", k, r.Header.Get(k)))
	}
	log.Printf("\n%s %s %s\nHost: %s\n%s\n", r.Method, r.URL.String(), r.Proto, r.Host, b.String())
}

// LoggingHandler is an http.Handler which provide logging for HTTP requests by using the specified Logger.
type LoggingHandler struct {
	Logger Logger
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.log(r)
}
