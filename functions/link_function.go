package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type linkFunction struct {
	url Expression
}

func newLinkFunction(args []Expression) (*linkFunction, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'link' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	return &linkFunction{url: args[0]}, nil
}

func (f *linkFunction) evaluate(vars map[string]interface{}, req *http.Request) (*HttpRsp, error) {
	a, err := f.url.Evaluate(vars, req)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	b, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	u, err := url.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %v", err)
	}
	rsp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %v", err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %v", err)
	}
	r := &HttpRsp{Body: string(body), Headers: rsp.Header, StatusCode: rsp.StatusCode}
	return r, nil
}
