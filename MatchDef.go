package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type MatchDef struct {
	Pattern  string       `json:"pattern"`
	Response *MatchedResp `json:"response"`
}

type MatchedResp struct {
	Body       string                 `json:"body"`
	Headers    map[string]interface{} `json:"headers"`
	StatusCode int                    `json:"status_code"`
}

func ParseMatchDef(j []byte) ([]*MatchDef, error) {
	var r []*MatchDef
	err := json.Unmarshal(j, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rsp *MatchedResp) ParseBody() (string, error) {
	body := rsp.Body
	start := strings.Index(body, "${")
	if start == 0 {
		end := len(body) - 1
		if body[end] != '}' {
			return "", fmt.Errorf("unexpected token '%c' at position '%d': expected '}'", body[end], end)
		}
		rest := body[2:end]
		start = strings.Index(rest, "(")
		if start <= 0 {
			return "", fmt.Errorf("expected token '('")
		}
		name := rest[0:start]
		end = len(rest) - 1
		if rest[end] != ')' {
			return "", fmt.Errorf("unexpected token '%c' at position '%d': expected ')'", rest[end], end)
		}
		arg := rest[start+1 : end]
		switch name {
		case "text":
			return arg, nil
		case "file":
			content, err := ioutil.ReadFile(arg)
			if err != nil {
				return "", err
			}
			return string(content), nil
		default:
			return "", fmt.Errorf("function '%s' is not supported", name)
		}
		return arg, nil
	}
	return body, nil
}
