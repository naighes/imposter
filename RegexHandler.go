package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type RegexHandler struct {
	routes []*regexRoute
}

type regexRoute struct {
	pattern *regexp.Regexp
	handler http.Handler
}

func (handler *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range handler.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (handler *RegexHandler) handleFunc(pattern *regexp.Regexp, h func(http.ResponseWriter, *http.Request)) {
	handler.routes = append(handler.routes, &regexRoute{pattern, http.HandlerFunc(h)})
}

func NewRegexHandler(defs []*MatchDef) (*RegexHandler, error) {
	r := RegexHandler{}
	for _, def := range defs {
		reg, err := regexp.Compile(def.Pattern)
		if err != nil {
			return nil, err
		}
		def := def
		// NOTE: should we manage the whole handler in a lazy manner?
		body, err := def.Response.ParseBody()
		if err != nil {
			return nil, err
		}
		r.handleFunc(reg, func(w http.ResponseWriter, r *http.Request) {
			for k, v := range def.Response.Headers {
				w.Header().Set(k, v.(string))
			}
			w.WriteHeader(def.Response.StatusCode)
			fmt.Fprintln(w, body)
		})
	}
	return &r, nil
}
