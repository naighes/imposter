package functions

import (
	"fmt"
	"net/http"
	"net/url"
)

type redirectFunction struct {
	url Expression
}

func newRedirectFunction(args []Expression) (*redirectFunction, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'redirect' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	return &redirectFunction{url: args[0]}, nil
}

func (f *redirectFunction) evaluate(vars map[string]interface{}, req *http.Request) (*HttpRsp, error) {
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
	r := &HttpRsp{Headers: h, StatusCode: 301}
	return r, nil
}
