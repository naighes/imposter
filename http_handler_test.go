package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naighes/imposter/functions"
)

type fakeExpression struct {
	rsp *functions.HttpRsp
}

type errorExpression struct {
	err string
}

func (e fakeExpression) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return e.rsp, nil
}

func (e errorExpression) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return nil, fmt.Errorf(e.err)
}

func TestFuncHttpHandlerNoErrors(t *testing.T) {
	const expectedStatusCode = 200
	const expectedBody = "some content"
	r := httptest.NewRecorder()
	p := func(string) (functions.Expression, error) {
		e := &fakeExpression{rsp: &functions.HttpRsp{StatusCode: expectedStatusCode, Body: expectedBody}}
		return e, nil
	}
	h := FuncHttpHandler{Content: "unrelevant content"}
	f, err := h.HandleFunc(p)
	if err != nil {
		t.Errorf("HandleFunc raised an error")
		return
	}
	f(r, nil)
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
		return
	}
	rsp := r.Result()
	var body []byte
	if body, err = ioutil.ReadAll(rsp.Body); err != nil {
		t.Errorf("cannot read body")
		return
	}
	if b := string(body); b != expectedBody {
		t.Errorf("expected body '%s'; got '%s'", expectedBody, b)
		return
	}
}

func TestFuncHttpHandlerWithErrors(t *testing.T) {
	const expectedStatusCode = 500
	r := httptest.NewRecorder()
	p := func(string) (functions.Expression, error) {
		e := &errorExpression{err: "some error"}
		return e, nil
	}
	h := FuncHttpHandler{Content: "unrelevant content"}
	f, err := h.HandleFunc(p)
	if err != nil {
		t.Errorf("HandleFunc raised an error")
		return
	}
	f(r, nil)
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
		return
	}
}

func TestFuncHttpHandlerWithoutHttpRsp(t *testing.T) {
	const expectedStatusCode = 500
	r := httptest.NewRecorder()
	p := func(string) (functions.Expression, error) {
		e := &functions.StringIdentity{Value: "some value"}
		return e, nil
	}
	h := FuncHttpHandler{Content: "unrelevant content"}
	f, err := h.HandleFunc(p)
	if err != nil {
		t.Errorf("HandleFunc raised an error")
		return
	}
	f(r, nil)
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
		return
	}
}
