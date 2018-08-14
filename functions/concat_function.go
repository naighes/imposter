package functions

import (
	"bytes"
	"fmt"
)

type concatFunction struct {
	args []Expression
}

func newConcatFunction(args []Expression) (Expression, error) {
	if l := len(args); l < 2 {
		return nil, fmt.Errorf("function 'concat' is expecting at least two arguments of type 'string'; found %d argument(s) instead", l)
	}
	r := concatFunction{args: args}
	return r, nil
}

func (f concatFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	var r bytes.Buffer
	for _, arg := range f.args {
		a, err := g(arg)
		if err != nil {
			return false, err
		}
		b, ok := a.(string)
		if !ok {
			r.WriteString(fmt.Sprintf("%v", a))
		} else {
			r.WriteString(b)
		}
	}
	return r.String(), nil
}

func (f concatFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}

func (f concatFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g)
}
