package functions

import (
	"fmt"
	"net/http"
)

type inFunction struct {
	source Expression
	item   Expression
}

func newInFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'in' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	r := inFunction{source: args[0], item: args[1]}
	return r, nil
}

func (f inFunction) evaluate(g func(Expression) (interface{}, error)) (interface{}, error) {
	a, err := g(f.source)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	var ok bool
	var left []interface{}
	if left, ok = a.([]interface{}); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'array'", a)
	}
	b, err := g(f.item)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	for _, el := range left {
		if el == b {
			return true, nil
		}
	}
	return false, nil
}

func (f inFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}

func (f inFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g)
}
