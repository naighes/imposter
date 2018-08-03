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

func newLinkFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'link' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	r := linkFunction{url: args[0]}
	return r, nil
}

func (f linkFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	a, err := f.url.Evaluate(ctx)
	if err != nil {
		return nil, err
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
	r := &HTTPRsp{Body: string(body), Headers: rsp.Header, StatusCode: rsp.StatusCode}
	return r, nil
}

func (f linkFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	a, err := f.url.Test(ctx)
	if err != nil {
		return nil, err
	}
	_, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	h := make(http.Header)
	r := &HTTPRsp{Body: "", Headers: h, StatusCode: 200}
	return r, nil
}
