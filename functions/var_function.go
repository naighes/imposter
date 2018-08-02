package functions

import (
	"fmt"
	"net/http"
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

func (f varFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(vars, req)
		}
	}(vars, req)
	return f.evaluate(g, vars)
}

func (f varFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	g := func(vars map[string]interface{}, req *http.Request) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(vars, req)
		}
	}(vars, req)
	return f.evaluate(g, vars)
}
