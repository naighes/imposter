package main

import (
	"fmt"
	"io/ioutil"
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

func parseDef(rsp *MatchRsp) func(http.ResponseWriter, *http.Request) {
	parseBody := rsp.ParseBody()
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := parseBody()
		if err != nil {
			writeError(w, err)
			return
		}
		if rsp.Headers != nil {
			for k, v := range rsp.Headers {
				w.Header().Set(k, v.(string))
			}
		}
		w.WriteHeader(rsp.StatusCode)
		fmt.Fprintf(w, body)
	}
}

func parseFunc(str string) (func(http.ResponseWriter, *http.Request), error) {
	name, arg, err := ParseFunc(str)
	if err != nil {
		return nil, err
	}
	switch name {
	case "link":
		return func(w http.ResponseWriter, r *http.Request) {
			rsp, err := http.Get(arg)
			defer rsp.Body.Close()
			if err != nil {
				writeError(w, err)
				return
			}
			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				writeError(w, err)
				return
			}
			for k, _ := range rsp.Header {
				w.Header().Set(k, rsp.Header.Get(k))
			}
			w.WriteHeader(rsp.StatusCode)
			fmt.Fprintf(w, string(body))
		}, nil
	case "redirect":
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", arg)
			w.WriteHeader(301)
		}, nil
	default:
		return nil, fmt.Errorf("function '%s' is not supported", name)
	}
}

func HandleFunc(o interface{}, options *ConfigOptions) (func(http.ResponseWriter, *http.Request), error) {
	var rsp MatchRsp
	err := mapstructure.Decode(o, &rsp)
	if err == nil {
		return parseDef(&rsp), nil
	}
	str, ok := o.(string)
	if ok {
		return parseFunc(str)
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
