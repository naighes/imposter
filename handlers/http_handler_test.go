package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naighes/imposter/cfg"
	"github.com/naighes/imposter/functions"
)

type fakeExpression struct {
	rsp *functions.HTTPRsp
}

func (e fakeExpression) Evaluate(ctx *functions.EvaluationContext) (interface{}, error) {
	return e.rsp, nil
}

func (e fakeExpression) Test(ctx *functions.EvaluationContext) (interface{}, error) {
	return e.rsp, nil
}

type errorExpression struct {
	err string
}

func (e errorExpression) Evaluate(ctx *functions.EvaluationContext) (interface{}, error) {
	return nil, fmt.Errorf(e.err)
}

func (e errorExpression) Test(ctx *functions.EvaluationContext) (interface{}, error) {
	return nil, fmt.Errorf(e.err)
}

type stringExpression struct {
	value string
}

func (e stringExpression) Evaluate(ctx *functions.EvaluationContext) (interface{}, error) {
	return e.value, nil
}

func (e stringExpression) Test(ctx *functions.EvaluationContext) (interface{}, error) {
	return e.value, nil
}

func TestFuncHTTPHandlerNoErrors(t *testing.T) {
	const expectedStatusCode = 200
	const expectedBody = "some content"
	r := httptest.NewRecorder()
	p := func(string) (functions.Expression, error) {
		e := &fakeExpression{rsp: &functions.HTTPRsp{StatusCode: expectedStatusCode, Body: expectedBody}}
		return e, nil
	}
	h := funcHTTPHandler{content: "unrelevant content"}
	f, err := h.handleFunc(p)
	if err != nil {
		t.Errorf("handleFunc raised an error")
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

func TestFuncHTTPHandlerWithErrors(t *testing.T) {
	const expectedStatusCode = 500
	r := httptest.NewRecorder()
	p := func(string) (functions.Expression, error) {
		e := &errorExpression{err: "some error"}
		return e, nil
	}
	h := funcHTTPHandler{content: "unrelevant content"}
	f, err := h.handleFunc(p)
	if err != nil {
		t.Errorf("handleFunc raised an error")
		return
	}
	f(r, nil)
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
		return
	}
}

func TestFuncHTTPHandlerWithoutHTTPRsp(t *testing.T) {
	const expectedStatusCode = 500
	r := httptest.NewRecorder()
	p := func(string) (functions.Expression, error) {
		e := &stringExpression{value: "some value"}
		return e, nil
	}
	h := funcHTTPHandler{content: "unrelevant content"}
	f, err := h.handleFunc(p)
	if err != nil {
		t.Errorf("handleFunc raised an error")
		return
	}
	f(r, nil)
	if r.Code != expectedStatusCode {
		t.Errorf("expected status code %d; got %d", expectedStatusCode, r.Code)
		return
	}
}

func TestNoBodyWhenMethodIsHead(t *testing.T) {
	const expectedBody = ""
	r := httptest.NewRecorder()
	c := cfg.MatchRsp{StatusCode: "${200}", Body: "some content"}
	h := matchRspHTTPHandler{content: &c}
	f, err := h.handleFunc(functions.ParseExpression)
	if err != nil {
		t.Errorf("handleFunc raised an error")
		return
	}
	f(r, &http.Request{Method: "HEAD"})
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

func TestHTTPCookie(t *testing.T) {
	r := httptest.NewRecorder()
	cookie := cfg.HTTPCookie{Value: "some value"}
	c := cfg.MatchRsp{StatusCode: "${200}", Body: "some content", Cookies: map[string]cfg.HTTPCookie{
		"cookie": cookie,
	}}
	h := matchRspHTTPHandler{content: &c}
	f, err := h.handleFunc(functions.ParseExpression)
	if err != nil {
		t.Errorf("handleFunc raised an error")
		return
	}
	f(r, &http.Request{Method: "HEAD"})
	cookies := r.Result().Cookies()
	if l := len(cookies); l != 1 {
		t.Errorf("expected '%d' cookie(s); got '%d' instead ", 1, l)
		return
	}
}
