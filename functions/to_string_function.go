package functions

import (
	"fmt"
	"strconv"
)

type toStringFunction struct {
	arg Expression
}

func newToStringFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'to_string' is expecting one argument; found %d argument(s) instead", l)
	}
	r := toStringFunction{arg: args[0]}
	return r, nil
}

func (f toStringFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	a, err := g(f.arg)
	if err != nil {
		return "", err
	}
	switch a.(type) {
	case string:
		return a, nil
	case int:
		return strconv.Itoa(a.(int)), nil
	case float64:
		return strconv.FormatFloat(a.(float64), 'E', -1, 64), nil
	case bool:
		return strconv.FormatBool(a.(bool)), nil
	default:
		return fmt.Sprintf("%v", a), nil
	}
}

func (f toStringFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}

func (f toStringFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}
