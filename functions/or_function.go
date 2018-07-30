package functions

import (
	"fmt"
	"net/http"
)

type orFunction struct {
	args []Expression
}

func newOrFunction(args []Expression) (*orFunction, error) {
	if l := len(args); l < 2 {
		return nil, fmt.Errorf("function 'or' is expecting at least two arguments of type 'bool'; found %d argument(s) instead", l)
	}
	return &orFunction{args: args}, nil
}

func (f *orFunction) evaluate(vars map[string]interface{}, req *http.Request) (bool, error) {
	r := false
	for _, arg := range f.args {
		a, err := arg.Evaluate(vars, req)
		if err != nil {
			return false, fmt.Errorf("%v", err)
		}
		b, ok := a.(bool)
		if !ok {
			return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", a)
		}
		r = r || b
	}
	return r, nil
}
