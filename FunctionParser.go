package main

import (
	"fmt"
	"strings"
)

func ParseFunc(str string) (string, string, error) {
	start := strings.Index(str, "${")
	if start != 0 {
		return "", "", fmt.Errorf("unexpected token '%c' at position '0': expected '${'", str[1])
	}
	end := len(str) - 1
	if str[end] != '}' {
		return "", "", fmt.Errorf("unexpected token '%c' at position '%d': expected '}'", str[end], end)
	}
	rest := str[2:end]
	start = strings.Index(rest, "(")
	if start <= 0 {
		return "", "", fmt.Errorf("expected token '('")
	}
	name := rest[0:start]
	end = len(rest) - 1
	if rest[end] != ')' {
		return "", "", fmt.Errorf("unexpected token '%c' at position '%d': expected ')'", rest[end], end)
	}
	arg := rest[start+1 : end]
	return name, arg, nil
}
