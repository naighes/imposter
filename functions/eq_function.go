package functions

import (
	"fmt"
	"net/http"
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
		return false, fmt.Errorf("%v", err)
	}
	b, err := g(f.right)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	return a == b, nil
}

func (f eqFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}

func (f eqFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}
