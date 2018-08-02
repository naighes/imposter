package functions

import (
	"fmt"
	"net/http"
)

type notFunction struct {
	arg Expression
}

func newNotFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'not' is expecting one argument of type 'bool'; found %d argument(s) instead", l)
	}
	r := notFunction{arg: args[0]}
	return r, nil
}

func (f notFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	a, err := g(f.arg)
	if err != nil {
		return false, err
	}
	b, ok := a.(bool)
	if !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", a)
	}
	return !b, nil
}

func (f notFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}

func (f notFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}
