package functions

import (
	"fmt"
	"net/http"
)

type httpHeaderFunction struct {
	name Expression
}

func newHttpHeaderFunction(args []Expression) (*httpHeaderFunction, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'http_header' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	return &httpHeaderFunction{name: args[0]}, nil
}

func (f *httpHeaderFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := f.name.Evaluate(vars, req)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if h := req.Header; h != nil {
		return req.Header.Get(b), nil
	} else {
		return "", fmt.Errorf("no http headers in source HTTP request")
	}
}
