package functions

import (
	"fmt"
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

func (f requestHTTPHostFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	return ctx.Req.Host, nil
}

func (f requestHTTPHostFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	return f.Evaluate(ctx)
}
