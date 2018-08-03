package functions

import (
	"fmt"
)

type requestURLFunction struct {
}

func newRequestURLFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 0 {
		return nil, fmt.Errorf("function 'request_url' is expecting no arguments; found %d argument(s) instead", l)
	}
	r := requestURLFunction{}
	return r, nil
}

func (f requestURLFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	if ctx.Req.URL == nil {
		return "", nil
	}
	return ctx.Req.URL.String(), nil
}

func (f requestURLFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	return f.Evaluate(ctx)
}
