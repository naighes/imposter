package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpHandler interface {
	HandleFunc() (func(http.ResponseWriter, *http.Request), error)
}

type FuncHttpHandler struct {
	Content string
	HttpGet func(string) (*http.Response, error)
}

func (h FuncHttpHandler) HandleFunc() (func(http.ResponseWriter, *http.Request), error) {
	name, arg, err := ParseFunc(h.Content)
	if err != nil {
		return nil, err
	}
	switch name {
	case "link":
		return func(w http.ResponseWriter, r *http.Request) {
			rsp, err := h.HttpGet(arg)
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

type MatchRspHttpHandler struct {
	Content *MatchRsp
}

func (h MatchRspHttpHandler) HandleFunc() (func(http.ResponseWriter, *http.Request), error) {
	rsp := h.Content
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
	}, nil
}
