package functions

import (
	"fmt"
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

func (f requestURLPathFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	if ctx.Req.URL == nil {
		return "", nil
	}
	return ctx.Req.URL.Path, nil
}

func (f requestURLPathFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	return f.Evaluate(ctx)
}
