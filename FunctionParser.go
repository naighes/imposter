package main

import (
	"fmt"
	"strings"
)

func ParseFunc(str string) (string, []string, error) {
	start := strings.Index(str, "${")
	if start != 0 {
		return "", nil, fmt.Errorf("unexpected token '%c' at position '0': expected '${'", str[1])
	}
	end := len(str) - 1
	if str[end] != '}' {
		return "", nil, fmt.Errorf("unexpected token '%c' at position '%d': expected '}'", str[end], end)
	}
	rest := str[2:end]
	start = strings.Index(rest, "(")
	if start <= 0 {
		return "", nil, fmt.Errorf("expected token '('")
	}
	name := rest[0:start]
	end = len(rest) - 1
	if rest[end] != ')' {
		return "", nil, fmt.Errorf("unexpected token '%c' at position '%d': expected ')'", rest[end], end)
	}
	arg := rest[start+1 : end]
	args, err := ParseArgs(arg)
	if err != nil {
		return "", nil, fmt.Errorf("an error has occurred on arguments parsing: %s", err)
	}
	return name, args, nil
}

func parseString(str string) (string, int, error) {
	start := 1
	for {
		s := str[start:]
		i := strings.Index(s, "\"")
		if i == -1 {
			return "", -1, fmt.Errorf("expected token '\"'")
		}
		if s[i-1] != '\\' {
			r := str[1 : start+i]
			i = i + start + 1
			return r, i, nil
		} else {
			start = start + i + 1
		}
	}
}

func parseNumber(str string) (string, int, error) {
	var chars []byte
	i := 0
	for {
		if i >= len(str) {
			return string(chars), i, nil
		}
		c := str[i]
		if c >= '0' && c <= '9' {
			chars = append(chars, c)
		} else if c == ',' {
			return string(chars), i + 1, nil
		} else if c != ' ' {
			return "", -1, fmt.Errorf("expected token ','; got %s", c)
		}
		i = i + 1
	}
}

func ParseArgs(str string) ([]string, error) {
	var r []string
	for {
		i := 0
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			break
		}
		var parser func(string) (string, int, error)
		var err error
		var token string
		switch str[i] {
		case '"':
			parser = parseString
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			parser = parseNumber
		default:
			return nil, fmt.Errorf("could not find a parser for the current token")
		}
		token, i, err = parser(str[i:])
		if err != nil {
			return nil, err
		}
		r = append(r, token)
		if i == len(str) {
			break
		}
		str = strings.TrimSpace(str[i:])
		if str[0] != ',' {
			return nil, fmt.Errorf("expected token ','; got '%s'", str[0])
		}
		str = str[1:]
	}
	return r, nil
}
