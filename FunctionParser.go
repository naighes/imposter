package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type expression interface {
	evaluate(map[string]interface{}) (interface{}, error)
}

type ifElse struct {
	guard expression
	left  expression
	right expression
}

type function struct {
	name string
	args []expression
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

type booleanIdentity struct {
	value string
}

func (e stringIdentity) evaluate(vars map[string]interface{}) (interface{}, error) {
	return e.value, nil
}

func (e integerIdentity) evaluate(vars map[string]interface{}) (interface{}, error) {
	return strconv.Atoi(e.value)
}

func (e floatIdentity) evaluate(vars map[string]interface{}) (interface{}, error) {
	return strconv.ParseFloat(e.value, 64)
}

func (e booleanIdentity) evaluate(vars map[string]interface{}) (interface{}, error) {
	return strconv.ParseBool(e.value)
}

func (e function) evaluate(vars map[string]interface{}) (interface{}, error) {
	switch e.name {
	case "link":
		return evaluateLink(e.args, vars)
	case "redirect":
		return evaluateRedirect(e.args, vars)
	case "file":
		return evaluateFile(e.args, vars)
	case "var":
		return evaluateVar(e.args, vars)
	case "and":
		return evaluateAnd(e.args, vars)
	case "or":
		return evaluateOr(e.args, vars)
	default:
		return nil, fmt.Errorf("function '%s' is not implemented", e.name)
	}
}

func (e ifElse) evaluate(vars map[string]interface{}) (interface{}, error) {
	guard, err := e.guard.evaluate(vars)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %s", err)
	}
	guardValue, ok := guard.(bool)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to bool", guard)
	}
	var exp expression
	if guardValue {
		exp = e.left
	} else {
		exp = e.right
	}
	val, err := exp.evaluate(vars)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %s", err)
	}
	return val, nil
}

type HttpRsp struct {
	Body       string
	Headers    http.Header
	StatusCode int
}

func evaluateLink(args []expression, vars map[string]interface{}) (interface{}, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'link' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars)
	if err != nil {
		return 0, fmt.Errorf("evaluation error: %s", err)
	}
	b, ok := a.(string)
	if !ok {
		return 0, fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	rsp, err := http.Get(b)
	defer rsp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("evaluation error: %s", err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return 0, fmt.Errorf("evaluation error: %s", err)
	}
	r := &HttpRsp{Body: string(body), Headers: rsp.Header, StatusCode: rsp.StatusCode}
	return r, nil
}

func evaluateRedirect(args []expression, vars map[string]interface{}) (interface{}, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'redirect' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %s", err)
	}
	b, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	h := make(http.Header)
	h.Set("Location", b)
	r := &HttpRsp{Headers: h, StatusCode: 301}
	return r, nil
}

func evaluateFile(args []expression, vars map[string]interface{}) (string, error) {
	if l := len(args); l != 1 {
		return "", fmt.Errorf("function 'file' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars)
	if err != nil {
		return "", fmt.Errorf("evaluation error: %s", err)
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	content, err := ioutil.ReadFile(b)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func evaluateVar(args []expression, vars map[string]interface{}) (string, error) {
	if l := len(args); l != 1 {
		return "", fmt.Errorf("function 'var' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars)
	if err != nil {
		return "", fmt.Errorf("evaluation error: %s", err)
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	if v, ok := vars[b]; ok {
		return v.(string), nil
	}
	return "", fmt.Errorf("evaluation error: cannot find a variable named '%s'", b)
}

func evaluateAnd(args []expression, vars map[string]interface{}) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'and' is expecting at least two arguments of type 'boolean'; found %d argument(s) instead", l)
	}
	r := true
	for _, arg := range args {
		a, err := arg.evaluate(vars)
		if err != nil {
			return false, fmt.Errorf("evaluation error: %s", err)
		}
		b, ok := a.(bool)
		if !ok {
			return false, fmt.Errorf("evaluation error: cannot convert value '%v' to boolean", a)
		}
		r = r && b
	}
	return r, nil
}

func evaluateOr(args []expression, vars map[string]interface{}) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'or' is expecting at least two arguments of type 'boolean'; found %d argument(s) instead", l)
	}
	r := false
	for _, arg := range args {
		a, err := arg.evaluate(vars)
		if err != nil {
			return false, fmt.Errorf("evaluation error: %s", err)
		}
		b, ok := a.(bool)
		if !ok {
			return false, fmt.Errorf("evaluation error: cannot convert value '%v' to boolean", a)
		}
		r = r || b
	}
	return r, nil
}

type parser func(string, int) (expression, int, error)

func stringParser(str string, start int) (expression, int, error) {
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

func numberParser(str string, start int) (expression, int, error) {
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
	var e expression
	if dots == 0 {
		e = &integerIdentity{value: str[start:end]}
	} else {
		e = &floatIdentity{value: str[start:end]}
	}
	return e, end, nil
}

func functionParser(str string, start int) (expression, int, error) {
	for {
		if start >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, start)
		}
		c := str[start]
		if c == ' ' {
			start = start + 1
			continue
		}
		if isLetterOrNumber(c) || c == '_' {
			break
		} else {
			return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d: expected ')'", c, start), str, start)
		}
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
	for {
		if end >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, end)
		}
		c := str[end]
		if c == ' ' {
			end = end + 1
			continue
		}
		if c == '(' {
			break
		}
		return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d", c, end), str, end)
	}
	var args []expression
	var err error
	args, end, err = parseArgs(str, end+1)
	if err != nil {
		return nil, -1, err
	}
	e := &function{name: name, args: args}
	return e, end, nil
}

func ifParser(str string, start int) (expression, int, error) {
	start = start + 2
	for {
		if start >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, start)
		}
		c := str[start]
		if c == ' ' {
			start = start + 1
			continue
		}
		if c == '(' {
			break
		} else {
			return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d: expected '('", c, start), str, start)
		}
	}
	// guard parsing
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
	for {
		if start >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, start)
		}
		c := str[start]
		if c == ' ' {
			start = start + 1
			continue
		}
		if c == ')' {
			break
		} else {
			return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d: expected ')'", c, start), str, start)
		}
	}
	// true condition expression parsing
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
	for {
		if start >= len(str) {
			return nil, -1, prettyError("unexpected end of string", str, start)
		}
		c := str[start]
		if c == ' ' {
			start = start + 1
			continue
		}
		if len(str) >= start+4 && str[start:start+4] == "else" {
			start = start + 4
			break
		} else {
			return nil, -1, prettyError(fmt.Sprintf("unexpected token '%c' at position %d: expected 'else'", c, start), str, start)
		}
	}
	// false condition expression parsing
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
		if c == ' ' {
			start = start + 1
			continue
		}
		if isLetter(c) {
			if len(str) >= start+4 && str[start:start+4] == "true" {
				return func(string, int) (expression, int, error) {
					e := &booleanIdentity{value: "true"}
					return e, start + 4, nil
				}, start, nil
			}
			if len(str) >= start+5 && str[start:start+5] == "false" {
				return func(string, int) (expression, int, error) {
					e := &booleanIdentity{value: "false"}
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
		if c == ')' {
			return nil, start, nil
		}
		return nil, -1, prettyError("could not find a parser for the current token", str, start)
	}
}

func parseArgs(str string, start int) ([]expression, int, error) {
	var r []expression
	hasToken := true
	for hasToken {
		var p parser
		var e expression
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
				return nil, -1, prettyError("expected token ')'", str, start)
			}
			c := str[start]
			start = start + 1
			if c == ')' {
				hasToken = false
				break
			}
			if c == ' ' {
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

func ParseExpression(str string) (expression, error) {
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
		e, start, err := p(str, start)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
	e := &stringIdentity{value: str}
	return e, nil
}

func prettyError(err string, str string, position int) error {
	const offset = 15
	l := int(math.Max(float64(position-offset), 0))
	r := int(math.Min(float64(position+offset), float64(len(str)-1)))
	s := str[l : r+1]
	var b bytes.Buffer
	b.WriteString(err)
	b.WriteString("\n")
	if l > 0 {
		b.WriteString("...")
		l = l - 3
	}
	b.WriteString(s)
	if r < len(str)-1 {
		b.WriteString("...")
	}
	b.WriteString("\n")
	for i := l; i < position; i++ {
		b.WriteString(" ")
	}
	b.WriteString("^")
	return fmt.Errorf(b.String())
}
