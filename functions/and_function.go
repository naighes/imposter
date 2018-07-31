package functions

import (
	"fmt"
	"net/http"
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
			return false, fmt.Errorf("%v", err)
		}
		b, ok := a.(bool)
		if !ok {
			return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", a)
		}
		r = r && b
	}
	return r, nil
}

func (f andFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}

func (f andFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}
