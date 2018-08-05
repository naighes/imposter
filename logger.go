package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type logger interface {
	log(r *http.Request)
}

type defaultLogger struct {
}

func (l *defaultLogger) log(r *http.Request) {
	if r == nil {
		return
	}
	var b bytes.Buffer
	for k := range r.Header {
		b.WriteString(fmt.Sprintf("%s: %s\n", k, r.Header.Get(k)))
	}
	log.Printf("\n%s %s %s\nHost: %s\n%s\n", r.Method, r.URL.String(), r.Proto, r.Host, b.String())
}

type loggingHandler struct {
	logger logger
	next   http.Handler
}

func (h *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.log(r)
	if h.next != nil {
		h.next.ServeHTTP(w, r)
	}
}
