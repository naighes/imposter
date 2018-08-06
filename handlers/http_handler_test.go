package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

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
		e := &functions.StringIdentity{Value: "some value"}
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
