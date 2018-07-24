package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func evaluateLink(args []expression, vars map[string]interface{}, req *http.Request) (interface{}, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'link' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
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

func evaluateRedirect(args []expression, vars map[string]interface{}, req *http.Request) (interface{}, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'redirect' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
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

func evaluateFile(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 1 {
		return "", fmt.Errorf("function 'file' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
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

func evaluateVar(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 1 {
		return "", fmt.Errorf("function 'var' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
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

func evaluateAnd(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'and' is expecting at least two arguments of type 'boolean'; found %d argument(s) instead", l)
	}
	r := true
	for _, arg := range args {
		a, err := arg.evaluate(vars, req)
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

func evaluateOr(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'or' is expecting at least two arguments of type 'boolean'; found %d argument(s) instead", l)
	}
	r := false
	for _, arg := range args {
		a, err := arg.evaluate(vars, req)
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

func evaluateNot(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l != 1 {
		return false, fmt.Errorf("function 'not' is expecting one argument of type 'boolean'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	b, ok := a.(bool)
	if !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to boolean", a)
	}
	return !b, nil
}

func evaluateHttpHeader(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 1 {
		return "", fmt.Errorf("function 'http_header' is expecting one argument of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
	if err != nil {
		return "", fmt.Errorf("evaluation error: %s", err)
	}
	b, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	if h := req.Header; h != nil {
		return req.Header.Get(b), nil
	} else {
		return "", fmt.Errorf("no http headers in source HTTP request")
	}
}

func evaluateEq(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'eq' is expecting two arguments; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	b, err := args[1].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	return a == b, nil
}

func evaluateNe(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'eq' is expecting two arguments; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	b, err := args[1].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	return a != b, nil
}

func evaluateContains(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'contains' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	var ok bool
	var left string
	if left, ok = a.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	b, err := args[1].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	var right string
	if right, ok = b.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	return strings.Index(left, right) != -1, nil
}

func evaluateRequestURL(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 0 {
		return "", fmt.Errorf("function 'request_url' is expecting no arguments; found %d argument(s) instead", l)
	}
	return req.URL.String(), nil
}

func evaluateRegexMatch(args []expression, vars map[string]interface{}, req *http.Request) (bool, error) {
	if l := len(args); l < 2 {
		return false, fmt.Errorf("function 'regex_match' is expecting two arguments of type 'string'; found %d argument(s) instead", l)
	}
	a, err := args[0].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	var ok bool
	var left string
	if left, ok = a.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	b, err := args[1].evaluate(vars, req)
	if err != nil {
		return false, fmt.Errorf("evaluation error: %s", err)
	}
	var right string
	if right, ok = b.(string); !ok {
		return false, fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
	}
	reg, err := regexp.Compile(right)
	if err != nil {
		return false, err
	}
	return reg.MatchString(left), nil
}

func evaluateRequestURLPath(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 0 {
		return "", fmt.Errorf("function 'request_url_path' is expecting no arguments; found %d argument(s) instead", l)
	}
	return req.URL.Path, nil
}

func evaluateRequestURLQuery(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	l := len(args)
	switch l {
	case 0:
		return req.URL.RawQuery, nil
	case 1:
		a, err := args[0].evaluate(vars, req)
		if err != nil {
			return "", fmt.Errorf("evaluation error: %s", err)
		}
		arg, ok := a.(string)
		if !ok {
			return "", fmt.Errorf("evaluation error: cannot convert value '%v' to string", a)
		}
		return req.URL.Query().Get(arg), nil
	default:
		return "", fmt.Errorf("function 'request_url_query' is expecting one or no arguments; found %d argument(s) instead", l)
	}
}

func evaluateRequestHTTPMethod(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 0 {
		return "", fmt.Errorf("function 'request_http_method' is expecting no arguments; found %d argument(s) instead", l)
	}
	return req.Method, nil
}

func evaluateRequestHost(args []expression, vars map[string]interface{}, req *http.Request) (string, error) {
	if l := len(args); l != 0 {
		return "", fmt.Errorf("function 'request_host' is expecting no arguments; found %d argument(s) instead", l)
	}
	return req.Host, nil
}
