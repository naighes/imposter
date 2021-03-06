package functions

import (
	"fmt"
	"strings"
)

type containsFunction struct {
	source Expression
	value  Expression
}

func newContainsFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'contains' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	r := containsFunction{source: args[0], value: args[1]}
	return r, nil
}

func (f containsFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	a, err := g(f.source)
	if err != nil {
		return false, err
	}
	var ok bool
	var left string
	if left, ok = a.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	b, err := g(f.value)
	if err != nil {
		return false, err
	}
	var right string
	if right, ok = b.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	return strings.Index(left, right) != -1, nil
}

func (f containsFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}

func (f containsFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}
