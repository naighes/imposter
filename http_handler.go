package main

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/naighes/imposter/functions"
	"github.com/spf13/cast"
)

type HttpHandler interface {
	HandleFunc(parse func(string) (functions.Expression, error)) (func(http.ResponseWriter, *http.Request), error)
}

type FuncHttpHandler struct {
	Content string
	Vars    map[string]interface{}
}

func (h FuncHttpHandler) HandleFunc(parse func(string) (functions.Expression, error)) (func(http.ResponseWriter, *http.Request), error) {
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
		rsp, ok := a.(*functions.HttpRsp)
		if !ok {
			writeError(w, fmt.Errorf("full response computing requires a function returning '*HttpRsp' (e.g. 'link', 'redirect', ...); got '%s' instead", reflect.TypeOf(a)))
			return
		}
		for k, _ := range rsp.Headers {
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

type MatchRspHttpHandler struct {
	Content *MatchRsp
	Vars    map[string]interface{}
}

func evaluateToString(e functions.Expression, vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := e.Evaluate(vars, req)
	if err != nil {
		return "", err
	}
	if b, ok := a.(string); ok {
		return b, nil
	}
	if b, err := cast.ToStringE(a); err == nil {
		return b, nil
	}
	return interfaceToString(a)
}

func (h MatchRspHttpHandler) HandleFunc(parse func(string) (functions.Expression, error)) (func(http.ResponseWriter, *http.Request), error) {
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
		if rsp.StatusCode > 0 {
			w.WriteHeader(rsp.StatusCode)
		} else {
			w.WriteHeader(200)
		}
		fmt.Fprintf(w, b)
	}, nil
}

// TODO: due to the lack of toString method in golang
func interfaceToString(i interface{}) (string, error) {
	a, ok := i.([]interface{})
	if !ok {
		return "", nil // TODO
	}
	var r []string
	for _, v := range a {
		e, err := cast.ToStringE(v)
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		r = append(r, e)
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(strings.Join(r, ","))
	buffer.WriteString("]")
	return buffer.String(), nil
}
