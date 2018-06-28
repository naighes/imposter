package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

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
	parseFunc := func(str string) (func(http.ResponseWriter, *http.Request), error) {
		start := strings.Index(str, "${")
		if start != 0 {
			return nil, fmt.Errorf("unexpected token '%c' at position '0': expected '${'", str[1])
		}
		end := len(str) - 1
		if str[end] != '}' {
			return nil, fmt.Errorf("unexpected token '%c' at position '%d': expected '}'", str[end], end)
		}
		rest := str[2:end]
		start = strings.Index(rest, "(")
		if start <= 0 {
			return nil, fmt.Errorf("expected token '('")
		}
		name := rest[0:start]
		end = len(rest) - 1
		if rest[end] != ')' {
			return nil, fmt.Errorf("unexpected token '%c' at position '%d': expected ')'", rest[end], end)
		}
		arg := rest[start+1 : end]
		switch name {
		case "link":
			return func(w http.ResponseWriter, r *http.Request) {
				rsp, err := http.Get(arg)
				if err != nil {
					// TODO: must be tested when 500
					w.WriteHeader(500)
					fmt.Fprintln(w, err)
					return
				}
				defer rsp.Body.Close()
				body, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					// TODO: must be tested when 500
					w.WriteHeader(500)
					fmt.Fprintln(w, err)
					return
				}
				w.WriteHeader(rsp.StatusCode)
				for k, _ := range rsp.Header {
					w.Header().Set(k, rsp.Header.Get(k))
				}
				fmt.Fprintln(w, string(body))
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
		return parseFunc(str)
	}
	// TODO: check for "link"
	return func(w http.ResponseWriter, r *http.Request) {
	}, nil
}
