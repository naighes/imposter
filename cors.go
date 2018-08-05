package main

import (
	"net/http"
)

type corsHandler struct {
	next    http.Handler
	enabled bool
}

func (h *corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.enabled {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	h.next.ServeHTTP(w, r)
}
