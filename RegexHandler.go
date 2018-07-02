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
	for _, def := range defs {
		reg, err := regexp.Compile(def.Pattern)
		if err != nil {
			return nil, err
		}
		f, err := HandleFunc(def.Response, config.Options)
		if err != nil {
			return nil, err
		}
		r.addRoute(reg, f)
	}
	return &r, nil
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "text/plain charset=utf-8")
	fmt.Fprintln(w, err)
}

func HandleFunc(o interface{}, options *ConfigOptions) (func(http.ResponseWriter, *http.Request), error) {
	parseDef := func(rsp *MatchRsp) func(http.ResponseWriter, *http.Request) {
		parseBody := rsp.ParseBody()
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(rsp.StatusCode)
			body, err := parseBody()
			if err != nil {
				writeError(w, err)
				return
			}
			if rsp.Headers != nil {
				for k, v := range rsp.Headers {
					fmt.Printf("%s: %s\n", k, v)
					w.Header().Set(k, v.(string))
				}
			}
			fmt.Fprintln(w, body)
		}
	}
	parseFunc := func(str string) (func(http.ResponseWriter, *http.Request), error) {
		name, arg, err := ParseFunc(str)
		if err != nil {
			return nil, err
		}
		switch name {
		case "link":
			return func(w http.ResponseWriter, r *http.Request) {
				rsp, err := http.Get(arg)
				if err != nil {
					writeError(w, err)
					return
				}
				defer rsp.Body.Close()
				body, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					writeError(w, err)
					return
				}
				w.WriteHeader(rsp.StatusCode)
				for k, _ := range rsp.Header {
					w.Header().Set(k, rsp.Header.Get(k))
				}
				fmt.Fprintln(w, string(body))
			}, nil
		case "redirect":
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(301)
				w.Header().Set("Location", arg)
			}, nil
		default:
			return nil, fmt.Errorf("function '%s' is not supported", name)
		}
	}
	var rsp MatchRsp
	err := mapstructure.Decode(o, &rsp)
	if err == nil {
		return parseDef(&rsp), nil
	}
	str, ok := o.(string)
	if ok {
		return enrichHeaders(parseFunc(str))
	}
	return nil, fmt.Errorf("operation is not supported")
}

func enrichHeaders(f func(http.ResponseWriter, *http.Request), err error) (func(http.ResponseWriter, *http.Request), error) {
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}, nil
}
