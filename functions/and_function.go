package functions

import (
	"fmt"
	"net/http"
)

type andFunction struct {
	args []Expression
}

func newAndFunction(args []Expression) (*andFunction, error) {
	if l := len(args); l < 2 {
		return nil, fmt.Errorf("function 'and' is expecting at least two arguments of type 'bool'; found %d argument(s) instead", l)
	}
	return &andFunction{args: args}, nil
}

func (f *andFunction) evaluate(vars map[string]interface{}, req *http.Request) (bool, error) {
	r := true
	for _, arg := range f.args {
		a, err := arg.Evaluate(vars, req)
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
