package functions

import (
	"fmt"
)

type eqFunction struct {
	left  Expression
	right Expression
}

func newEqFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'eq' is expecting two arguments; found %d argument(s) instead", l)
	}
	r := eqFunction{left: args[0], right: args[1]}
	return r, nil
}

func (f eqFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	a, err := g(f.left)
	if err != nil {
		return false, err
	}
	b, err := g(f.right)
	if err != nil {
		return false, err
	}
	return a == b, nil
}

func (f eqFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}

func (f eqFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}
