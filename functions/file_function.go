package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type fileFunction struct {
	path Expression
}

func newFileFunction(args []Expression) (*fileFunction, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'file' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	return &fileFunction{path: args[0]}, nil
}

func (f *fileFunction) evaluate(vars map[string]interface{}, req *http.Request) (string, error) {
	a, err := f.path.Evaluate(vars, req)
	if err != nil {
		return "", fmt.Errorf("%v", err)
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
