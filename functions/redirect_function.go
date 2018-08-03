package functions

import (
	"fmt"
	"net/http"
	"net/url"
)

type redirectFunction struct {
	url        Expression
	statusCode Expression
}

func newRedirectFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'redirect' is expecting two arguments ('string', 'int'); found %d argument(s) instead", l)
	}
	r := redirectFunction{url: args[0], statusCode: args[1]}
	return r, nil
}

func (f redirectFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	a, err := f.url.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	b, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	c, err := f.statusCode.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	statusCode, ok := c.(int)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'int'", c)
	}
	u, err := url.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %v", err)
	}
	if statusCode < 300 || statusCode >= 400 {
		return nil, fmt.Errorf("evaluation error: expected status code '3XX'; got '%d' instead", statusCode)
	}
	h := make(http.Header)
	h.Set("Location", u.String())
	r := &HTTPRsp{Headers: h, StatusCode: statusCode}
	return r, nil
}

func (f redirectFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	a, err := f.url.Test(ctx)
	if err != nil {
		return nil, err
	}
	b, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	h := make(http.Header)
	h.Set("Location", b)
	r := &HTTPRsp{Headers: h, StatusCode: 301}
	return r, nil
}
