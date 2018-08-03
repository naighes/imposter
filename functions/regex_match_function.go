package functions

import (
	"fmt"
	"regexp"
)

type regexMatchFunction struct {
	source  Expression
	pattern Expression
}

func newRegexMatchFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 2 {
		return nil, fmt.Errorf("function 'regex_match' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	r := regexMatchFunction{source: args[0], pattern: args[1]}
	return r, nil
}

func (f regexMatchFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	a, err := f.source.Evaluate(ctx)
	if err != nil {
		return false, err
	}
	var ok bool
	var left string
	if left, ok = a.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	b, err := f.pattern.Evaluate(ctx)
	if err != nil {
		return false, err
	}
	var right string
	if right, ok = b.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	reg, err := regexp.Compile(right)
	if err != nil {
		return false, err
	}
	return reg.MatchString(left), nil
}

func (f regexMatchFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	a, err := f.source.Test(ctx)
	if err != nil {
		return false, err
	}
	var ok bool
	var left string
	if left, ok = a.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	b, err := f.pattern.Test(ctx)
	if err != nil {
		return false, err
	}
	var right string
	if right, ok = b.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to 'string'", a)
	}
	reg, err := regexp.Compile(right)
	if err != nil {
		return false, nil
	}
	return reg.MatchString(left), nil
}
