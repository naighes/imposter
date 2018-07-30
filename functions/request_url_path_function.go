package functions

import (
	"fmt"
	"net/http"
)

type requestURLPathFunction struct {
}

func newRequestURLPathFunction(args []Expression) (*requestURLPathFunction, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_url_path' is expecting no arguments; found %d argument(s) instead", l)
	}
	return &requestURLPathFunction{}, nil
}

func (f *requestURLPathFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	if req.URL == nil {
		return "", nil
	}
	return req.URL.Path, nil
}
