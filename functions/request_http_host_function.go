package functions

import (
	"fmt"
	"net/http"
)

type requestHTTPHostFunction struct {
}

func newRequestHTTPHostFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_http_host' is expecting no arguments; found %d argument(s) instead", l)
	}
	r := requestHTTPHostFunction{}
	return r, nil
}

func (f requestHTTPHostFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return req.Host, nil
}

func (f requestHTTPHostFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return f.Evaluate(vars, req)
}
