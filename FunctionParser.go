package main

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

type expression interface {
	evaluate(map[string]interface{}, *http.Request) (interface{}, error)
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

type arrayIdentity struct {
	elements []expression
}

func (e stringIdentity) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return e.value, nil
}

func (e integerIdentity) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return strconv.Atoi(e.value)
}

func (e floatIdentity) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return strconv.ParseFloat(e.value, 64)
}

func (e booleanIdentity) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	return strconv.ParseBool(e.value)
}

func (e arrayIdentity) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	var r []interface{}
	var t reflect.Type
	for index, element := range e.elements {
		a, err := element.evaluate(vars, req)
		if err != nil {
			return nil, fmt.Errorf("evaluation error: %s", err)
		}
		if index > 0 && t != reflect.TypeOf(a) {
			return nil, fmt.Errorf("evaluation error: mixed type arrays are not allowed")
		}
		t = reflect.TypeOf(a)
		r = append(r, a)
	}

	return r, nil
}

func (e function) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	switch e.name {
	case "link":
		return evaluateLink(e.args, vars, req)
	case "redirect":
		return evaluateRedirect(e.args, vars, req)
	case "file":
		return evaluateFile(e.args, vars, req)
	case "var":
		return evaluateVar(e.args, vars, req)
	case "and":
		return evaluateAnd(e.args, vars, req)
	case "or":
		return evaluateOr(e.args, vars, req)
	case "not":
		return evaluateNot(e.args, vars, req)
	case "http_header":
		return evaluateHttpHeader(e.args, vars, req)
	case "eq":
		return evaluateEq(e.args, vars, req)
	case "ne":
		return evaluateNe(e.args, vars, req)
	case "contains":
		return evaluateContains(e.args, vars, req)
	case "request_url":
		return evaluateRequestURL(e.args, vars, req)
	case "request_url_path":
		return evaluateRequestURLPath(e.args, vars, req)
	case "request_url_query":
		return evaluateRequestURLQuery(e.args, vars, req)
	case "regex_match":
		return evaluateRegexMatch(e.args, vars, req)
	default:
		return nil, fmt.Errorf("function '%s' is not implemented", e.name)
	}
}

func (e ifElse) evaluate(vars map[string]interface{}, req *http.Request) (interface{}, error) {
	guard, err := e.guard.evaluate(vars, req)
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
	val, err := exp.evaluate(vars, req)
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

func functionParser(str string, start int) (expression, int, error) {
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
	var args []expression
	args, end, err = parseArgs(str, end+1, ')')
	if err != nil {
		return nil, -1, err
	}
	e := &function{name: name, args: args}
	return e, end, nil
}

func ifParser(str string, start int) (expression, int, error) {
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

func arrayParser(str string, start int) (expression, int, error) {
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
		if c == '[' {
			return arrayParser, start, nil
		}
		if c == ')' {
			return nil, start, nil
		}
		return nil, -1, prettyError("could not find a parser for the current token", str, start)
	}
}

func parseArgs(str string, start int, endToken byte) ([]expression, int, error) {
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
		if start < position && end > position {
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
