package functions

import (
	"fmt"
	"net/http"
	"strings"
)

type containsFunction struct {
	source Expression
	value  Expression
}

func newContainsFunction(args []Expression) (*containsFunction, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'contains' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	return &containsFunction{source: args[0], value: args[1]}, nil
}

func (f *containsFunction) evaluate(vars map[string]interface{}, req *http.Request) (bool, error) {
	a, err := f.source.Evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	var ok bool
	var left string
	if left, ok = a.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	b, err := f.value.Evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	var right string
	if right, ok = b.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	return strings.Index(left, right) != -1, nil
}
