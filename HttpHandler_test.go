package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBla(t *testing.T) {
	r := httptest.NewRecorder()
	get := func(string) (*http.Response, error) {
		return (&ResponseMock{}).MakeResponse(), fmt.Errorf("raised error")
	}
	h := FuncHttpHandler{Content: "${link(http://fak.eurl)}", HttpGet: get}
	f, err := h.HandleFunc()
	if err != nil {
		t.Errorf("HandleFunc raised an error")
	}
	f(r, nil)
	const expected = 500
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}
