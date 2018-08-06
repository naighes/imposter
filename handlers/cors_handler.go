package handlers

import (
	"net/http"
)

type CorsHandler struct {
	Enabled bool
}

func (h *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.Enabled {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
