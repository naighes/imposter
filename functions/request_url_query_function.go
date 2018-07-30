package functions

import (
	"fmt"
	"net/http"
)

type requestURLQueryFunction struct {
	arg Expression
}

func newRequestURLQueryFunction(args []Expression) (*requestURLQueryFunction, error) {
	l := len(args)
	switch l {
	case 0:
		return &requestURLQueryFunction{}, nil
	case 1:
		return &requestURLQueryFunction{arg: args[0]}, nil
	default:
		return nil, fmt.Errorf("function 'request_url_query' is expecting one or no arguments; found %d argument(s) instead", l)
	}
}

func evaluateWithArgument(arg Expression, vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := arg.Evaluate(vars, req)
	if err != nil {
		return "", fmt.Errorf("%v", err)
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

func (f *requestURLQueryFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	if f.arg == nil {
		return evaluateWithoutArgument(vars, req)
	}
	return evaluateWithArgument(f.arg, vars, req)
}
