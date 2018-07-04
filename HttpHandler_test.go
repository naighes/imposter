package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLinkWhenHttpGetRaisesAnError(t *testing.T) {
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
	const expectedStatusCode = 500
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
	}
}

func TestLinkWhenHttpGetReturnsOk(t *testing.T) {
	const expectedStatusCode = 200
	const expectedBody = "Hello Test!"
	r := httptest.NewRecorder()
	body := []byte(expectedBody)
	get := func(string) (*http.Response, error) {
		return (&ResponseMock{StatusCode: expectedStatusCode, Body: body}).MakeResponse(), nil
	}
	h := FuncHttpHandler{Content: "${link(http://fak.eurl)}", HttpGet: get}
	f, err := h.HandleFunc()
	if err != nil {
		t.Errorf("HandleFunc raised an error")
	}
	f(r, nil)
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
	}
	rsp := r.Result()
	var content []byte
	if content, err = ioutil.ReadAll(rsp.Body); err != nil {
		t.Errorf("cannot read body")
	}
	if c := string(content); c != expectedBody {
		t.Errorf("expected body %s; got %s", expectedBody, c)
	}
}

func TestRedirect(t *testing.T) {
	const expectedUrl = "http://fak.eurl"
	r := httptest.NewRecorder()
	h := FuncHttpHandler{Content: fmt.Sprintf("${redirect(%s)}", expectedUrl)}
	f, err := h.HandleFunc()
	if err != nil {
		t.Errorf("HandleFunc raised an error")
	}
	f(r, nil)
	const expectedStatusCode = 301
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
	}
	location := r.Header().Get("Location")
	if location != expectedUrl {
		t.Errorf("expected Location header %s; got %s", expectedUrl, location)
	}
}
