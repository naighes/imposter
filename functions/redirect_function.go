package functions

import (
	"fmt"
	"net/http"
	"net/url"
)

type redirectFunction struct {
	url Expression
}

func newRedirectFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'redirect' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	r := redirectFunction{url: args[0]}
	return r, nil
}

func (f redirectFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	a, err := f.url.Evaluate(vars, req)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	b, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	u, err := url.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %v", err)
	}
	h := make(http.Header)
	h.Set("Location", u.String())
	r := &HTTPRsp{Headers: h, StatusCode: 301}
	return r, nil
}

func (f redirectFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	a, err := f.url.Test(vars, req)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
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
