package functions

import (
	"fmt"
	"net/http"
)

type requestURLFunction struct {
}

func newRequestURLFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_url' is expecting no arguments; found %d argument(s) instead", l)
	}
	r := requestURLFunction{}
	return r, nil
}

func (f requestURLFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	if req.URL == nil {
		return "", nil
	}
	return req.URL.String(), nil
}

func (f requestURLFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return f.Evaluate(vars, req)
}
