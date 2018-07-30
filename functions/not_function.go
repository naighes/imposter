package functions

import (
	"fmt"
	"net/http"
)

type notFunction struct {
	arg Expression
}

func newNotFunction(args []Expression) (*notFunction, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'not' is expecting one argument of type 'bool'; found %d argument(s) instead", l)
	}
	return &notFunction{arg: args[0]}, nil
}

func (f *notFunction) evaluate(vars map[string]interface{}, req *http.Request) (bool, error) {
	a, err := f.arg.Evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	b, ok := a.(bool)
	if !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", a)
	}
	return !b, nil
}
