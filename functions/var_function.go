package functions

import (
	"fmt"
	"net/http"
)

type varFunction struct {
	name Expression
}

func newVarFunction(args []Expression) (*varFunction, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'var' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	return &varFunction{name: args[0]}, nil
}

func (f *varFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := f.name.Evaluate(vars, req)
	if err != nil {
		return "", fmt.Errorf("%v", err)
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
