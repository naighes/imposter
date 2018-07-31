package functions

import (
	"fmt"
	"net/http"
)

type requestURLPathFunction struct {
}

func newRequestURLPathFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_url_path' is expecting no arguments; found %d argument(s) instead", l)
	}
	r := requestURLPathFunction{}
	return r, nil
}

func (f requestURLPathFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	if req.URL == nil {
		return "", nil
	}
	return req.URL.Path, nil
}

func (f requestURLPathFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return f.Evaluate(vars, req)
}
