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
}

func (h FuncHttpHandler) HandleFunc(parse func(string) (expression, error)) (func(http.ResponseWriter, *http.Request), error) {
	e, err := parse(h.Content)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := e.evaluate()
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
}

func (h MatchRspHttpHandler) HandleFunc(parse func(string) (expression, error)) (func(http.ResponseWriter, *http.Request), error) {
	rsp := h.Content
	e, err := parse(rsp.Body)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		a, err := e.evaluate()
		if err != nil {
			writeError(w, err)
			return
		}
		b, ok := a.(string)
		if !ok {
			// TODO: better error message
			writeError(w, fmt.Errorf("expected string; got '%v' instead", a))
			return
		}
		if rsp.Headers != nil {
			// TODO: evaluate headers
			for k, v := range rsp.Headers {
				w.Header().Set(k, v.(string))
			}
		}
		w.WriteHeader(rsp.StatusCode)
		fmt.Fprintf(w, b)
	}, nil
}
