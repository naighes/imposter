package functions

import (
	"fmt"
	"net/http"
)

type requestURLQueryFunction struct {
	arg Expression
}

func newRequestURLQueryFunction(args []Expression) (Expression, error) {
	l := len(args)
	switch l {
	case 0:
		r := requestURLQueryFunction{}
		return r, nil
	case 1:
		r := requestURLQueryFunction{arg: args[0]}
		return r, nil
	default:
		return nil, fmt.Errorf("function 'request_url_query' is expecting one or no arguments; found %d argument(s) instead", l)
	}
}

func evaluateWithArgument(arg Expression, vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := arg.Evaluate(vars, req)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if req.URL == nil {
		return "", nil
	}
	return req.URL.Query().Get(b), nil
}

// TODO: remove duplication
func testWithArgument(arg Expression, vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := arg.Test(vars, req)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if req.URL == nil {
		return "", nil
	}
	return req.URL.Query().Get(b), nil
}

func evaluateWithoutArgument(vars map[string]interface{}, req *http.Request) (string, error) {
	if req.URL == nil {
		return "", nil
	}
	return req.URL.RawQuery, nil
}

func (f requestURLQueryFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	if f.arg == nil {
		return evaluateWithoutArgument(vars, req)
	}
	return evaluateWithArgument(f.arg, vars, req)
}

func (f requestURLQueryFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	if f.arg == nil {
		return evaluateWithoutArgument(vars, req)
	}
	return testWithArgument(f.arg, vars, req)
}
