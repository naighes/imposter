package functions

import (
	"fmt"
	"net/http"
)

type requestHTTPHeaderFunction struct {
	name Expression
}

func newRequestHTTPHeaderFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'request_http_header' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	r := requestHTTPHeaderFunction{name: args[0]}
	return r, nil
}

func (f requestHTTPHeaderFunction) evaluate(g func(Expression) (interface{}, error), req *http.Request) (interface{}, error) {
	a, err := g(f.name)
	if err != nil {
		return "", err
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

func (f requestHTTPHeaderFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g, req)
}

func (f requestHTTPHeaderFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g, req)
}
