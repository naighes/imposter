package functions

import (
	"fmt"
	"math/rand"
	"net/http"
)

type rndIntFunction struct {
	min Expression
	max Expression
}

func newRndIntFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'rnd_int' is expecting two arguments of type 'int'; found %d argument(s) instead", l)
	}
	r := rndIntFunction{min: args[0], max: args[1]}
	return r, nil
}

func (f rndIntFunction) evaluate(g func(Expression) (interface{}, error), req *http.Request) (interface{}, error) {
	a, err := g(f.min)
	if err != nil {
		return 0, err
	}
	b, ok := a.(int)
	if !ok {
		return 0, fmt.Errorf("evaluation error: cannot convert value '%v' to 'int'", a)
	}
	c, err := g(f.max)
	if err != nil {
		return 0, err
	}
	d, ok := c.(int)
	if !ok {
		return 0, fmt.Errorf("evaluation error: cannot convert value '%v' to 'int'", c)
	}
	if b > d {
		return 0, fmt.Errorf("evaluation error: 'min' argument must be lower than 'max'")
	}
	return randomInt(b, d), nil
}

func (f rndIntFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Req)
}

func (f rndIntFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Req)
}

func randomInt(min, max int) int {
	r := rand.New(src)
	return min + r.Intn(max-min)
}
