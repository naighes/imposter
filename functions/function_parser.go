package functions

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type EvaluationContext struct {
	Vars map[string]interface{}
	Req  *http.Request
}

type ExpressionParser = func(string) (Expression, error)

// Expression is an interface representing the result
// of a parsing operation.
type Expression interface {
	Evaluate(*EvaluationContext) (interface{}, error)
	Test(*EvaluationContext) (interface{}, error)
}

type ifElse struct {
	guard Expression
	left  Expression
	right Expression
}

type function struct {
	name string
	args []Expression
}

type stringIdentity struct {
	value string
}

type integerIdentity struct {
	value string
}

type floatIdentity struct {
	value string
}

type boolIdentity struct {
	value string
}

type arrayIdentity struct {
	elements []Expression
}

func (e stringIdentity) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	return e.value, nil
}

func (e stringIdentity) Test(ctx *EvaluationContext) (interface{}, error) {
	return e.Evaluate(ctx)
}

func (e integerIdentity) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	return strconv.Atoi(e.value)
}

func (e integerIdentity) Test(ctx *EvaluationContext) (interface{}, error) {
	return e.Evaluate(ctx)
}

func (e floatIdentity) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	return strconv.ParseFloat(e.value, 64)
}

func (e floatIdentity) Test(ctx *EvaluationContext) (interface{}, error) {
	return e.Evaluate(ctx)
}

func (e boolIdentity) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	return strconv.ParseBool(e.value)
}

func (e boolIdentity) Test(ctx *EvaluationContext) (interface{}, error) {
	return e.Evaluate(ctx)
}

func (e arrayIdentity) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	f := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return e.evaluate(f)
}

func (e arrayIdentity) Test(ctx *EvaluationContext) (interface{}, error) {
	f := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return e.evaluate(f)
}

func (e arrayIdentity) evaluate(f func(Expression) (interface{}, error)) (interface{}, error) {
	var r []interface{}
	var t reflect.Type
	for index, element := range e.elements {
		a, err := f(element)
		if err != nil {
			return nil, err
		}
		if index > 0 && t != reflect.TypeOf(a) {
			return nil, fmt.Errorf("evaluation error: mixed type arrays are not allowed")
		}
		switch a.(type) {
		case int, string, bool, float64:
			t = reflect.TypeOf(a)
			r = append(r, a)
		default:
			return nil, fmt.Errorf("array support is limited to 'int', 'string', 'bool', 'float64': found '%v' instead", t)
		}
	}

	return r, nil
}

type Evaluate func(*EvaluationContext) (interface{}, error)

func getEvaluationFunc(name string) (func(args []Expression) (Expression, error), error) {
	var b func(args []Expression) (Expression, error)
	switch name {
	case "link":
		b = newLinkFunction
	case "redirect":
		b = newRedirectFunction
	case "file":
		b = newFileFunction
	case "var":
		b = newVarFunction
	case "and":
		b = newAndFunction
	case "or":
		b = newOrFunction
	case "not":
		b = newNotFunction
	case "request_http_header":
		b = newRequestHTTPHeaderFunction
	case "eq":
		b = newEqFunction
	case "ne":
		b = newNeFunction
	case "contains":
		b = newContainsFunction
	case "request_url":
		b = newRequestURLFunction
	case "request_url_path":
		b = newRequestURLPathFunction
	case "request_url_query":
		b = newRequestURLQueryFunction
	case "request_http_method":
		b = newRequestHTTPMethodFunction
	case "request_http_host":
		b = newRequestHTTPHostFunction
	case "regex_match":
		b = newRegexMatchFunction
	case "in":
		b = newInFunction
	case "to_string":
		b = newToStringFunction
	case "rnd_string":
		b = newRndStringFunction
	default:
		return nil, fmt.Errorf("could not find a built-in function with name '%s'", name)
	}
	return b, nil
}

func (e function) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	b, err := getEvaluationFunc(e.name)
	if err != nil {
		return nil, err
	}
	f, err := b(e.args)
	if err == nil {
		return f.Evaluate(ctx)
	}
	return nil, err
}

func (e function) Test(ctx *EvaluationContext) (interface{}, error) {
	b, err := getEvaluationFunc(e.name)
	if err != nil {
		return nil, err
	}
	f, err := b(e.args)
	if err == nil {
		return f.Test(ctx)
	}
	return nil, err
}

func (e ifElse) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	guard, err := e.guard.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	guardValue, ok := guard.(bool)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", guard)
	}
	var exp Expression
	if guardValue {
		exp = e.left
	} else {
		exp = e.right
	}
	val, err := exp.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (e ifElse) Test(ctx *EvaluationContext) (interface{}, error) {
	guard, err := e.guard.Test(ctx)
	if err != nil {
		return nil, err
	}
	_, ok := guard.(bool)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to 'bool'", guard)
	}
	left, err := e.left.Test(ctx)
	if err != nil {
		return nil, err
	}
	right, err := e.right.Test(ctx)
	if err != nil {
		return nil, err
	}
	a := reflect.TypeOf(left)
	b := reflect.TypeOf(right)
	if a != b {
		return nil, fmt.Errorf("type mismatch: cannot convert type '%v' to '%v'", a, b)
	}
	return left, nil
}

// HTTPRsp represents an abstraction for an HTTP response message.
type HTTPRsp struct {
	Body       string
	Headers    http.Header
	StatusCode int
}

type parser func(string, int) (Expression, int, error)

func stringParser(str string, start int) (Expression, int, error) {
	end := start + 1
	for {
		if end >= len(str) {
			break
		}
		c := str[end]
		if c == '"' {
			if str[end-1] == '\\' {
				end = end + 1
				continue
			} else {
				break
			}
		} else {
			end = end + 1
			continue
		}
	}
	e := &stringIdentity{value: str[start+1 : end]}
	return e, end + 1, nil
}

func numberParser(str string, start int) (Expression, int, error) {
	end := start
	dots := 0
	for {
		if end >= len(str) {
			break
		}
		c := str[end]
		if isNumber(c) {
			end = end + 1
			continue
		}
		if c == '.' {
			if dots == 0 {
				end = end + 1
				dots = dots + 1
				continue
			}
			return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d: expected ')'", c, end), str, end)
		}
		break
	}
	var e Expression
	if dots == 0 {
		e = &integerIdentity{value: str[start:end]}
	} else {
		e = &floatIdentity{value: str[start:end]}
	}
	return e, end, nil
}

func nextToken(str string, start int, match func(byte, int) bool) (int, error) {
	for {
		if start >= len(str) {
			return -1, prettyError("unexpected end of string", str, start)
		}
		c := str[start]
		if c == ' ' || c == '\n' || c == '\t' {
			start = start + 1
			continue
		}
		if match(c, start) {
			return start, nil
		}
		return -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d", c, start), str, start)
	}
}

func functionParser(str string, start int) (Expression, int, error) {
	start, err := nextToken(str, start, func(c byte, s int) bool {
		return isLetterOrNumber(c) || c == '_'
	})
	if err != nil {
		return nil, -1, err
	}
	end := start
	for {
		if end >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, end)
		}
		c := str[end]
		if isLetterOrNumber(c) || c == '_' {
			end = end + 1
			continue
		}
		break
	}
	if end == start {
		return nil, -1, prettyError("expected function name", str, end)
	}
	name := str[start:end]
	end, err = nextToken(str, end, func(c byte, s int) bool {
		return c == '('
	})
	if err != nil {
		return nil, -1, err
	}
	var args []Expression
	args, end, err = parseArgs(str, end+1, ')')
	if err != nil {
		return nil, -1, err
	}
	e := &function{name: name, args: args}
	return e, end, nil
}

func ifParser(str string, start int) (Expression, int, error) {
	var err error
	start = start + 2
	start, err = nextToken(str, start, func(c byte, s int) bool {
		return c == '('
	})
	if err != nil {
		return nil, -1, err
	}
	start = start + 1
	guardParser, start, err := getParser(str, start)
	if err != nil {
		return nil, -1, err
	}
	// TODO: ensure parser is not nil
	guard, start, err := guardParser(str, start)
	if err != nil {
		return nil, -1, err
	}
	start, err = nextToken(str, start, func(c byte, s int) bool {
		return c == ')'
	})
	if err != nil {
		return nil, -1, err
	}
	start = start + 1
	leftParser, start, err := getParser(str, start)
	if err != nil {
		return nil, -1, err
	}
	// TODO: ensure parser is not nil
	left, start, err := leftParser(str, start)
	if err != nil {
		return nil, -1, err
	}
	// check for else
	start, err = nextToken(str, start, func(c byte, s int) bool {
		return len(str) >= s+4 && str[s:s+4] == "else"
	})
	if err != nil {
		return nil, -1, err
	}
	start = start + 4
	rightParser, start, err := getParser(str, start)
	if err != nil {
		return nil, -1, err
	}
	// TODO: ensure parser is not nil
	right, start, err := rightParser(str, start)
	if err != nil {
		return nil, -1, err
	}
	e := &ifElse{guard: guard, left: left, right: right}
	return e, start, nil
}

func arrayParser(str string, start int) (Expression, int, error) {
	start = start + 1
	args, end, err := parseArgs(str, start, ']')
	if err != nil {
		return nil, -1, err
	}
	e := &arrayIdentity{elements: args}
	return e, end, nil
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isLetterOrNumber(c byte) bool {
	return isNumber(c) || isLetter(c)
}

func isNumber(c byte) bool {
	return (c >= '0' && c <= '9')
}

func getParser(str string, start int) (parser, int, error) {
	for {
		if start >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, start)
		}
		c := str[start]
		if c == ' ' || c == '\n' || c == '\t' {
			start = start + 1
			continue
		}
		if isLetter(c) {
			if len(str) >= start+4 && str[start:start+4] == "true" {
				return func(string, int) (Expression, int, error) {
					e := &boolIdentity{value: "true"}
					return e, start + 4, nil
				}, start, nil
			}
			if len(str) >= start+5 && str[start:start+5] == "false" {
				return func(string, int) (Expression, int, error) {
					e := &boolIdentity{value: "false"}
					return e, start + 5, nil
				}, start, nil
			}
			if len(str) >= start+2 && str[start:start+2] == "if" {
				return ifParser, start, nil
			}
			return functionParser, start, nil
		}
		if c == '"' {
			return stringParser, start, nil
		}
		if isNumber(c) {
			return numberParser, start, nil
		}
		if c == '[' {
			return arrayParser, start, nil
		}
		if c == ')' {
			return nil, start, nil
		}
		return nil, -1, prettyError("could not find a parser for the current token", str, start)
	}
}

func parseArgs(str string, start int, endToken byte) ([]Expression, int, error) {
	var r []Expression
	hasToken := true
	for hasToken {
		var p parser
		var e Expression
		var err error
		p, start, err = getParser(str, start)
		if err != nil {
			return nil, -1, err
		}
		if p != nil {
			e, start, err = p(str, start)
			if err != nil {
				return nil, -1, err
			}
			r = append(r, e)
		}
		for {
			if start >= len(str) {
				return nil, -1, prettyError("unexpected end of string: expected token ')'", str, start)
			}
			c := str[start]
			start = start + 1
			if c == endToken {
				hasToken = false
				break
			}
			if c == ' ' || c == '\n' || c == '\t' {
				continue
			}
			if c != ',' {
				return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d: expected ','", c, start-1), str, start-1)
			}
			break
		}
		if !hasToken {
			break
		}
	}
	return r, start, nil
}

// ParseExpression takes a string as an input and
// evaluates its content by returning an Expression.
func ParseExpression(str string) (Expression, error) {
	if strings.Index(str, "${") == 0 && str[len(str)-1] == '}' {
		str = str[2 : len(str)-1]
		start := 0
		p, start, err := getParser(str, start)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, prettyError("could not find a parser for the current token", str, start)
		}
		e, _, err := p(str, start)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
	e := &stringIdentity{value: str}
	return e, nil
}

func prettyError(err string, str string, position int) error {
	const offset = 30
	l := int(math.Max(float64(position-offset), 0))
	r := int(math.Min(float64(position+offset), float64(len(str)-1)))
	s := str[l : r+1]
	var b bytes.Buffer
	b.WriteString(err)
	b.WriteString("\n")
	if l > 0 {
		b.WriteString("...\n")
	}
	scanner := bufio.NewScanner(strings.NewReader(s))
	acc := l
	for scanner.Scan() {
		token := scanner.Text()
		b.WriteString(token)
		b.WriteString("\n")
		start := acc
		end := acc + len(token) + 1
		if start < position && end >= position {
			for i := acc; i < position; i++ {
				b.WriteString(" ")
			}
			b.WriteString("^")
			b.WriteString("\n")
		}
		acc = end
	}
	if r < len(str)-1 {
		b.WriteString("...")
	}
	return fmt.Errorf(b.String())
}
