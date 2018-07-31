package functions

import (
	"fmt"
	"net/http"
)

type requestHTTPMethodFunction struct {
}

func newRequestHTTPMethodFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_http_method' is expecting no arguments; found %d argument(s) instead", l)
	}
	r := requestHTTPMethodFunction{}
	return r, nil
}

func (f requestHTTPMethodFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return req.Method, nil
}

func (f requestHTTPMethodFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return f.Evaluate(vars, req)
}
