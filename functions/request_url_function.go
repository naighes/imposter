package functions

import (
	"fmt"
	"net/http"
)

type requestURLFunction struct {
}

func newRequestURLFunction(args []Expression) (*requestURLFunction, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_url' is expecting no arguments; found %d argument(s) instead", l)
	}
	return &requestURLFunction{}, nil
}

func (f *requestURLFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	if req.URL == nil {
		return "", nil
	}
	return req.URL.String(), nil
}
