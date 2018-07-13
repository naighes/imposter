package main

import (
	"fmt"
	"net/http"
)

type HttpHandler interface {
	HandleFunc() (func(http.ResponseWriter, *http.Request), error)
}

type FuncHttpHandler struct {
	Content string
}

func (h FuncHttpHandler) HandleFunc() (func(http.ResponseWriter, *http.Request), error) {
	e, err := ParseExpression(h.Content)
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
			// TODO: better error message
			writeError(w, fmt.Errorf("expected full HttpRsp; got '%v' instead", a))
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

func (h MatchRspHttpHandler) HandleFunc() (func(http.ResponseWriter, *http.Request), error) {
	rsp := h.Content
	e, err := ParseExpression(rsp.Body)
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
