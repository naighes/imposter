package functions

import (
	"fmt"
)

type orFunction struct {
	args []Expression
}

func newOrFunction(args []Expression) (Expression, error) {
	if l := len(args); l < 2 {
		return nil, fmt.Errorf("function 'or' is expecting at least two arguments of type 'bool'; found %d argument(s) instead", l)
	}
	r := orFunction{args: args}
	return r, nil
}

func (f orFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	r := false
	for _, arg := range f.args {
		a, err := g(arg)
		if err != nil {
			return false, err
		}
		b, ok := a.(bool)
		if !ok {
			return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", a)
		}
		r = r || b
	}
	return r, nil
}

func (f orFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}

func (f orFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}
