package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/naighes/imposter/functions"
)

type HTTPHandler interface {
	HandleFunc(parse func(string) (functions.Expression, error)) (func(http.ResponseWriter, *http.Request), error)
}

type FuncHTTPHandler struct {
	Content string
	Vars    map[string]interface{}
}

func (h FuncHTTPHandler) HandleFunc(parse func(string) (functions.Expression, error)) (func(http.ResponseWriter, *http.Request), error) {
	e, err := parse(h.Content)
	if err != nil {
		return nil, err
	}
	vars := h.Vars
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := e.Evaluate(vars, r)
		if err != nil {
			writeError(w, err)
			return
		}
		rsp, ok := a.(*functions.HTTPRsp)
		if !ok {
			writeError(w, fmt.Errorf("full response computing requires a function returning '*HTTPRsp' (e.g. 'link', 'redirect', ...); got '%s' instead", reflect.TypeOf(a)))
			return
		}
		for k := range rsp.Headers {
			w.Header().Set(k, rsp.Headers.Get(k))
		}
		if rsp.StatusCode > 0 {
			w.WriteHeader(rsp.StatusCode)
		} else {
			w.WriteHeader(200)
		}
		fmt.Fprintf(w, rsp.Body)
	}, nil
}

type MatchRspHTTPHandler struct {
	Content *MatchRsp
	Vars    map[string]interface{}
}

func (h MatchRspHTTPHandler) HandleFunc(parse func(string) (functions.Expression, error)) (func(http.ResponseWriter, *http.Request), error) {
	rsp := h.Content
	e, err := parse(rsp.Body)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]functions.Expression)
	if rsp.Headers != nil {
		for k, v := range rsp.Headers {
			header, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("expected a value of type 'string'; got '%s' instead", reflect.TypeOf(v))
			}
			he, err := parse(header)
			if err != nil {
				return nil, err
			}
			headers[k] = he
		}
	}
	vars := h.Vars
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := e.Evaluate(vars, r)
		if err != nil {
			writeError(w, err)
			return
		}
		for k, v := range headers {
			v1, err := v.Evaluate(vars, r)
			if err != nil {
				writeError(w, err)
				return
			}
			w.Header().Set(k, fmt.Sprintf("%v", v1))
		}
		if rsp.StatusCode > 0 {
			w.WriteHeader(rsp.StatusCode)
		} else {
			w.WriteHeader(200)
		}
		fmt.Fprintf(w, "%v", b)
	}, nil
}
