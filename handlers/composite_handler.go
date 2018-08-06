package handlers

import (
	"net/http"
)

type CompositeHandler struct {
	NestedHandlers []http.Handler
}

func (h *CompositeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range h.NestedHandlers {
		e.ServeHTTP(w, r)
	}
}
