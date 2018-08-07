package handlers

import (
	"net/http"
)

// CompositeHandler type executes a set of http.Handler sequentially.
type CompositeHandler struct {
	NestedHandlers []http.Handler
}

func (h *CompositeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range h.NestedHandlers {
		e.ServeHTTP(w, r)
	}
}
