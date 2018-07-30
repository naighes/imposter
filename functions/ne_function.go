package functions

import (
	"fmt"
	"net/http"
)

type neFunction struct {
	left  Expression
	right Expression
}

func newNeFunction(args []Expression) (*neFunction, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'ne' is expecting two arguments; found %d argument(s) instead", l)
	}
	return &neFunction{left: args[0], right: args[1]}, nil
}

func (f *neFunction) evaluate(vars map[string]interface{}, req *http.Request) (bool, error) {
	a, err := f.left.Evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	b, err := f.right.Evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	return a != b, nil
}
