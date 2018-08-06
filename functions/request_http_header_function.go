package functions

import (
	"fmt"
	"net/http"
)

type requestHTTPHeaderFunction struct {
	name Expression
}

func newRequestHTTPHeaderFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'request_http_header' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	r := requestHTTPHeaderFunction{name: args[0]}
	return r, nil
}

func (f requestHTTPHeaderFunction) evaluate(g func(Expression) (interface{}, error), req *http.Request) (interface{}, error) {
	a, err := g(f.name)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if h := req.Header; h != nil {
		return req.Header.Get(b), nil
	}
	return "", fmt.Errorf("no http headers in source HTTP request")
}

func (f requestHTTPHeaderFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Req)
}

func (f requestHTTPHeaderFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Req)
}
