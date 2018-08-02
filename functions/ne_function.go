package functions

import (
	"fmt"
	"net/http"
)

type neFunction struct {
	left  Expression
	right Expression
}

func newNeFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'ne' is expecting two arguments; found %d argument(s) instead", l)
	}
	r := neFunction{left: args[0], right: args[1]}
	return r, nil
}

func (f neFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}

func (f neFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}

func (f neFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	a, err := g(f.left)
	if err != nil {
		return false, err
	}
	b, err := g(f.right)
	if err != nil {
		return false, err
	}
	return a != b, nil
}
