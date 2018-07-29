package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type RegexHandler struct {
	routes []*regexRoute
	vars   map[string]interface{}
}

type regexRoute struct {
	expression expression
	method     string
	latency    time.Duration
	handler    http.Handler
}

func (handler *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range handler.routes {
		// TODO: X-Forwarded-Host?
		a, err := route.expression.evaluate(handler.vars, r)
		if err != nil {
			writeError(w, err)
			return
		}
		b, ok := a.(bool)
		if !ok {
			writeError(w, fmt.Errorf("rule_expression requires a 'bool' expression: found '%v' instead", reflect.TypeOf(a)))
			return
		}
		if b {
			if route.latency > 0 {
				time.Sleep(route.latency * time.Millisecond)
			}
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// TODO: not sure about just returning not found...
	http.NotFound(w, r)
}

func (handler *RegexHandler) addRoute(expression expression, method string, latency time.Duration, h func(http.ResponseWriter, *http.Request)) {
	handler.routes = append(handler.routes, &regexRoute{expression, method, latency, http.HandlerFunc(h)})
}

func NewRegexHandler(config *Config) (*RegexHandler, error) {
	defs := config.Defs
	var options *ConfigOptions
	if config.Options == nil {
		options = &ConfigOptions{}
	} else {
		options = config.Options
	}
	var vars map[string]interface{}
	if config.Vars == nil {
		vars = make(map[string]interface{})
	} else {
		vars = config.Vars
	}
	r := RegexHandler{}
	r.vars = vars
	for _, def := range defs {
		rule, err := ParseExpression(def.RuleExpression)
		if err != nil {
			return nil, err
		}
		f, err := HandleFunc(def.Response, options, vars)
		if err != nil {
			return nil, err
		}
		method, err := getMethod(def)
		if err != nil {
			return nil, err
		}
		if def.Latency < 0 {
			return nil, fmt.Errorf("latency requires a value greater than zero")
		}
		r.addRoute(rule, method, def.Latency, enrichHeaders(f, options))
	}
	return &r, nil
}

func getMethod(def *MatchDef) (string, error) {
	if def.Method == "" {
		return "*", nil
	}
	m := strings.ToUpper(def.Method)
	switch m {
	case "OPTIONS", "HEAD", "GET", "POST", "PUT", "DELETE", "TRACE", "*":
		return m, nil
	default:
		return "", fmt.Errorf("HTTP method '%s' is not supported", m)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain charset=utf-8")
	w.WriteHeader(500)
	fmt.Fprintf(w, err.Error())
}

func HandleFunc(o interface{}, options *ConfigOptions, vars map[string]interface{}) (func(http.ResponseWriter, *http.Request), error) {
	var rsp MatchRsp
	err := mapstructure.Decode(o, &rsp)
	if err == nil {
		return MatchRspHttpHandler{Content: &rsp, Vars: vars}.HandleFunc(ParseExpression)
	}
	str, ok := o.(string)
	if ok {
		return FuncHttpHandler{Content: str, Vars: vars}.HandleFunc(ParseExpression)
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
