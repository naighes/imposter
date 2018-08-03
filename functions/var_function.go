package functions

import (
	"fmt"
)

type varFunction struct {
	name Expression
}

func newVarFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'var' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	r := varFunction{name: args[0]}
	return r, nil
}

func (f varFunction) evaluate(g func(Expression) (interface{}, error), vars map[string]interface{}) (interface{}, error) {
	a, err := g(f.name)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	if v, ok := vars[b]; ok {
		return v.(string), nil
	}
	return "", fmt.Errorf("evaluation error: cannot find a variable named '%s'", b)
}

func (f varFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Vars)
}

func (f varFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Vars)
}
