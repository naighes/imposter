package main

import (
	"fmt"
	"net/http"
	"reflect"
)

type HttpHandler interface {
	HandleFunc(parse func(string) (expression, error)) (func(http.ResponseWriter, *http.Request), error)
}

type FuncHttpHandler struct {
	Content string
	Vars    map[string]interface{}
}

func (h FuncHttpHandler) HandleFunc(parse func(string) (expression, error)) (func(http.ResponseWriter, *http.Request), error) {
	e, err := parse(h.Content)
	if err != nil {
		return nil, err
	}
	vars := h.Vars
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := e.evaluate(vars, r)
		if err != nil {
			writeError(w, err)
			return
		}
		rsp, ok := a.(*HttpRsp)
		if !ok {
			writeError(w, fmt.Errorf("full response computing requires a function returning '*HttpRsp' (e.g. 'link', 'redirect', ...); got '%s' instead", reflect.TypeOf(a)))
			return
		}
		for k, _ := range rsp.Headers {
			w.Header().Set(k, rsp.Headers.Get(k))
		}
		w.WriteHeader(rsp.StatusCode)
		fmt.Fprintf(w, rsp.Body)
	}, nil
}

type MatchRspHttpHandler struct {
	Content *MatchRsp
	Vars    map[string]interface{}
}

func evaluateToString(e expression, vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := e.evaluate(vars, req)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("expected a value of type string; got '%s' instead", reflect.TypeOf(a))
	}
	return b, nil
}

func (h MatchRspHttpHandler) HandleFunc(parse func(string) (expression, error)) (func(http.ResponseWriter, *http.Request), error) {
	rsp := h.Content
	e, err := parse(rsp.Body)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]expression)
	if rsp.Headers != nil {
		for k, v := range rsp.Headers {
			header, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("expected a value of type string; got '%s' instead", reflect.TypeOf(v))
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
		b, err := evaluateToString(e, vars, r)
		if err != nil {
			writeError(w, err)
			return
		}
		for k, v := range headers {
			v1, err := evaluateToString(v, vars, r)
			if err != nil {
				writeError(w, err)
				return
			}
			w.Header().Set(k, v1)
		}
		w.WriteHeader(rsp.StatusCode)
		fmt.Fprintf(w, b)
	}, nil
}
