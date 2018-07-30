package functions

import (
	"fmt"
	"net/http"
)

type requestHTTPHostFunction struct {
}

func newRequestHTTPHostFunction(args []Expression) (*requestHTTPHostFunction, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_http_host' is expecting no arguments; found %d argument(s) instead", l)
	}
	return &requestHTTPHostFunction{}, nil
}

func (f *requestHTTPHostFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	return req.Host, nil
}
