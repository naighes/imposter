package functions

import (
	"fmt"
	"net/http"
)

type requestHTTPMethodFunction struct {
}

func newRequestHTTPMethodFunction(args []Expression) (*requestHTTPMethodFunction, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_http_method' is expecting no arguments; found %d argument(s) instead", l)
	}
	return &requestHTTPMethodFunction{}, nil
}

func (f *requestHTTPMethodFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	return req.Method, nil
}
