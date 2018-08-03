package functions

import (
	"fmt"
)

type requestURLQueryFunction struct {
	arg Expression
}

func newRequestURLQueryFunction(args []Expression) (Expression, error) {
	l := len(args)
	switch l {
	case 0:
		r := requestURLQueryFunction{}
		return r, nil
	case 1:
		r := requestURLQueryFunction{arg: args[0]}
		return r, nil
	default:
		return nil, fmt.Errorf("function 'request_url_query' is expecting one or no arguments; found %d argument(s) instead", l)
	}
}

func evaluateWithArgument(arg Expression, ctx *EvaluationContext) (string, error) {
	a, err := arg.Evaluate(ctx)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if ctx.Req.URL == nil {
		return "", nil
	}
	return ctx.Req.URL.Query().Get(b), nil
}

// TODO: remove duplication
func testWithArgument(arg Expression, ctx *EvaluationContext) (string, error) {
	a, err := arg.Test(ctx)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if ctx.Req.URL == nil {
		return "", nil
	}
	return ctx.Req.URL.Query().Get(b), nil
}

func evaluateWithoutArgument(ctx *EvaluationContext) (string, error) {
	if ctx.Req.URL == nil {
		return "", nil
	}
	return ctx.Req.URL.RawQuery, nil
}

func (f requestURLQueryFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	if f.arg == nil {
		return evaluateWithoutArgument(ctx)
	}
	return evaluateWithArgument(f.arg, ctx)
}

func (f requestURLQueryFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	if f.arg == nil {
		return evaluateWithoutArgument(ctx)
	}
	return testWithArgument(f.arg, ctx)
}
