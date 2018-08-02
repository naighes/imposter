package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/naighes/imposter/functions"
)

type HTTPHandler interface {
	HandleFunc(parse functions.ExpressionParser) (func(http.ResponseWriter, *http.Request), error)
}

type FuncHTTPHandler struct {
	Content string
	Vars    map[string]interface{}
}

func (h FuncHTTPHandler) HandleFunc(parse functions.ExpressionParser) (func(http.ResponseWriter, *http.Request), error) {
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
			writeError(w, fmt.Errorf("expected a positive 'int' value for status code; got '%d' instead", rsp.StatusCode))
		}
		fmt.Fprintf(w, rsp.Body)
	}, nil
}

type MatchRspHTTPHandler struct {
	Content *MatchRsp
	Vars    map[string]interface{}
}

func (h MatchRspHTTPHandler) HandleFunc(parse functions.ExpressionParser) (func(http.ResponseWriter, *http.Request), error) {
	rsp := h.Content
	e1, err := parse(rsp.Body)
	if err != nil {
		return nil, err
	}
	headers, err := rsp.parseHeaders(parse)
	if err != nil {
		return nil, err
	}
	e2, err := rsp.parseStatusCode(parse)
	if err != nil {
		return nil, err
	}
	vars := h.Vars
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := e1.Evaluate(vars, r)
		if err != nil {
			writeError(w, err)
			return
		}
		statusCode, err := evaluateStatusCode(e2, vars, r)
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
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "%v", b)
	}, nil
}

func evaluateStatusCode(e functions.Expression, vars map[string]interface{}, r *http.Request) (int, error) {
	s, err := e.Evaluate(vars, r)
	if err != nil {
		return 0, err
	}
	statusCode, ok := s.(int)
	if !ok {
		return 0, fmt.Errorf("expected an 'int' value for status code; got '%v' instead", reflect.TypeOf(s))
	}
	if statusCode <= 0 {
		return 0, fmt.Errorf("expected a positive 'int' value for status code; got '%d' instead", statusCode)
	}
	return statusCode, nil
}
