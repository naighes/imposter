package functions

import (
	"fmt"
	"net/http"
)

type inFunction struct {
	source Expression
	item   Expression
}

func newInFunction(args []Expression) (*inFunction, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'in' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	return &inFunction{source: args[0], item: args[1]}, nil
}

func (f *inFunction) evaluate(vars map[string]interface{}, req *http.Request) (bool, error) {
	a, err := f.source.Evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	var ok bool
	var left []interface{}
	if left, ok = a.([]interface{}); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'array'", a)
	}
	b, err := f.item.Evaluate(vars, req)
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
