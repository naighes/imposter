package functions

import (
	"fmt"
)

type andFunction struct {
	args []Expression
}

func newAndFunction(args []Expression) (Expression, error) {
	if l := len(args); l < 2 {
		return nil, fmt.Errorf("function 'and' is expecting at least two arguments of type 'bool'; found %d argument(s) instead", l)
	}
	r := andFunction{args: args}
	return r, nil
}

func (f andFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	r := true
	for _, arg := range f.args {
		a, err := g(arg)
		if err != nil {
			return false, err
		}
		b, ok := a.(bool)
		if !ok {
			return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", a)
		}
		r = r && b
	}
	return r, nil
}

func (f andFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}

func (f andFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}
