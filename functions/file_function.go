package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type fileFunction struct {
	path Expression
}

func newFileFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'file' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	r := fileFunction{path: args[0]}
	return r, nil
}

func (f fileFunction) Evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	a, err := f.path.Evaluate(vars, req)
	if err != nil {
		return "", err
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	content, err := ioutil.ReadFile(b)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (f fileFunction) Test(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	a, err := f.path.Test(vars, req)
	if err != nil {
		return "", err
	}
	_, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	return "", nil
}
