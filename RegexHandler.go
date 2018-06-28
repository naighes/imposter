package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/mitchellh/mapstructure"
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

func (handler *RegexHandler) addRoute(pattern *regexp.Regexp, h func(http.ResponseWriter, *http.Request)) {
	handler.routes = append(handler.routes, &regexRoute{pattern, http.HandlerFunc(h)})
}

func NewRegexHandler(defs []*MatchDef) (*RegexHandler, error) {
	r := RegexHandler{}
	for _, def := range defs {
		reg, err := regexp.Compile(def.Pattern)
		if err != nil {
			return nil, err
		}
		f, err := HandleFunc(def.Response)
		if err != nil {
			return nil, err
		}
		r.addRoute(reg, f)
	}
	return &r, nil
}

func HandleFunc(o interface{}) (func(http.ResponseWriter, *http.Request), error) {
	parseDef := func(rsp *MatchRsp) func(http.ResponseWriter, *http.Request) {
		parseBody := rsp.ParseBody()
		return func(w http.ResponseWriter, r *http.Request) {
			body, err := parseBody()
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintln(w, err)
				return
			}
			w.WriteHeader(rsp.StatusCode)
			for k, v := range rsp.Headers {
				w.Header().Set(k, v.(string))
			}
			fmt.Fprintln(w, body)
		}
	}
	var rsp MatchRsp
	err := mapstructure.Decode(o, &rsp)
	if err != nil {
		// TODO: check for "link"
		return func(w http.ResponseWriter, r *http.Request) {
		}, nil
	}
	return parseDef(&rsp), nil
}
