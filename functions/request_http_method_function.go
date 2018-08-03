package functions

import (
	"fmt"
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

func (f requestHTTPMethodFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	return ctx.Req.Method, nil
}

func (f requestHTTPMethodFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	return f.Evaluate(ctx)
}
