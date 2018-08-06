package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type Logger interface {
	log(r *http.Request)
}

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

type LoggingHandler struct {
	Logger Logger
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.log(r)
}
