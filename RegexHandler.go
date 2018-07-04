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
		// TODO: host and X-Forwarded-Host
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// TODO: not sure about just returning not found...
	http.NotFound(w, r)
}

func (handler *RegexHandler) addRoute(pattern *regexp.Regexp, h func(http.ResponseWriter, *http.Request)) {
	handler.routes = append(handler.routes, &regexRoute{pattern, http.HandlerFunc(h)})
}

func NewRegexHandler(config *Config) (*RegexHandler, error) {
	r := RegexHandler{}
	defs := config.Defs
	var options *ConfigOptions
	if config.Options == nil {
		options = &ConfigOptions{}
	} else {
		options = config.Options
	}
	for _, def := range defs {
		reg, err := regexp.Compile(def.Pattern)
		if err != nil {
			return nil, err
		}
		f, err := HandleFunc(def.Response, options)
		if err != nil {
			return nil, err
		}
		r.addRoute(reg, enrichHeaders(f, options))
	}
	return &r, nil
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain charset=utf-8")
	w.WriteHeader(500)
	fmt.Fprintf(w, err.Error())
}

func HandleFunc(o interface{}, options *ConfigOptions) (func(http.ResponseWriter, *http.Request), error) {
	var rsp MatchRsp
	err := mapstructure.Decode(o, &rsp)
	if err == nil {
		return MatchRspHttpHandler{Content: &rsp}.HandleFunc()
	}
	str, ok := o.(string)
	if ok {
		return FuncHttpHandler{Content: str, HttpGet: http.Get}.HandleFunc()
	}
	return nil, fmt.Errorf("operation is not supported")
}

func enrichHeaders(f func(http.ResponseWriter, *http.Request), options *ConfigOptions) func(http.ResponseWriter, *http.Request) {
	o := options
	return func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w, o)
		f(w, r)
	}
}

func setCorsHeaders(w http.ResponseWriter, options *ConfigOptions) {
	// NOTE: preflighted requests are not handled: we may came back on this in future
	if options.Cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
}
