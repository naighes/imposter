package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestIntegerIdentity(t *testing.T) {
	str := "${123}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	_, ok := token.(*IntegerIdentity)
	if !ok {
		t.Errorf("expected type '*IntegerIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
}

func TestFloatIdentity(t *testing.T) {
	str := "${1.123}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	_, ok := token.(*FloatIdentity)
	if !ok {
		t.Errorf("expected type '*FloatIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
}

func TestStringIdentity(t *testing.T) {
	str := "abc"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	_, ok := token.(*StringIdentity)
	if !ok {
		t.Errorf("expected type '*StringIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
}

func TestBoolIdentity(t *testing.T) {
	str := "${true}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	b, ok := token.(*BoolIdentity)
	if !ok {
		t.Errorf("expected type '*BoolIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	e, err := b.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Errorf("evaluation error")
		return
	}
	v, ok := e.(bool)
	if !ok {
		t.Errorf("expected type 'bool'; got '%s'", reflect.TypeOf(e))
		return
	}
	if !v {
		t.Errorf("expected 'bool' value 'true'; got '%v' instead", v)
		return
	}
}

func TestArrayIdentity(t *testing.T) {
	str := "${[123, 456]}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*ArrayIdentity)
	if !ok {
		t.Errorf("expected type '*ArrayIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	e, err := f.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Errorf("evaluation error: %v", err)
		return
	}
	v, ok := e.([]interface{})
	if !ok {
		t.Errorf("expected type 'map[interface{}]bool'; got '%s'", reflect.TypeOf(e))
		return
	}
	if l := len(v); l != 2 {
		t.Errorf("expected array of length 2; got '%d'", l)
		return
	}
}

func TestMixedTypeArrayIdentity(t *testing.T) {
	str := `${[123, "hello"]}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*ArrayIdentity)
	if !ok {
		t.Errorf("expected type '*ArrayIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	_, err = f.Evaluate(make(map[string]interface{}), &http.Request{})
	if err == nil {
		t.Errorf("expected error")
		return
	}
}

func TestOrEvaluation(t *testing.T) {
	str := "${or(false, true)}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	b, ok := token.(*Function)
	if !ok {
		t.Errorf("expected type '*Function'; got '%s'", reflect.TypeOf(token))
		return
	}
	e, err := b.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Errorf("evaluation error")
		return
	}
	v, ok := e.(bool)
	if !ok {
		t.Errorf("expected type 'bool'; got '%s'", reflect.TypeOf(e))
		return
	}
	if !v {
		t.Errorf("expected 'bool' value 'true'; got '%v' instead", v)
		return
	}
}

func TestAndEvaluation(t *testing.T) {
	str := `${
				and(true, 
					true)}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	b, ok := token.(*Function)
	if !ok {
		t.Errorf("expected type '*Function'; got '%s'", reflect.TypeOf(token))
		return
	}
	e, err := b.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Errorf("evaluation error")
		return
	}
	v, ok := e.(bool)
	if !ok {
		t.Errorf("expected type 'bool'; got '%s'", reflect.TypeOf(e))
		return
	}
	if !v {
		t.Errorf("expected 'bool' value 'true'; got '%v' instead", v)
		return
	}
}

func TestBlockWithJustIdentity(t *testing.T) {
	str := `${"abc"}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	_, ok := token.(*StringIdentity)
	if !ok {
		t.Errorf("expected type '*StringIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
}

func TestEmptyString(t *testing.T) {
	str := ""
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	si, ok := token.(*StringIdentity)
	if !ok {
		t.Errorf("expected type '*StringIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	s, err := si.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Errorf("evaluation error")
		return
	}
	v, ok := s.(string)
	if !ok {
		t.Errorf("expected type 'string'; got '%s'", reflect.TypeOf(s))
		return
	}
	if v != "" {
		t.Errorf("expected an empty string; got '%s'", v)
		return
	}
}

func TestFunctionWithArguments(t *testing.T) {
	str := `${
				f(
					"12345" , 
					987
				) 
			}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*Function)
	if !ok {
		t.Errorf("expected type '*Function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.Name != "f" {
		t.Errorf("expected Function named '%s'; got '%s'", "f", f.Name)
		return
	}
	if l := len(f.Args); l != 2 {
		t.Errorf("expected '%d' argument(s); got '%d'", 2, l)
		return
	}
	arg1, ok := f.Args[0].(*StringIdentity)
	if !ok {
		t.Errorf("expected argument of type '*StringIdentity'; got '%s'", reflect.TypeOf(arg1))
		return
	}
	arg2, ok := f.Args[1].(*IntegerIdentity)
	if !ok {
		t.Errorf("expected argument of type '*IntegerIdentity'; got '%s'", reflect.TypeOf(arg2))
		return
	}
}

func TestFunctionWithoutArguments(t *testing.T) {
	str := "${  f  (     )  }"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*Function)
	if !ok {
		t.Errorf("expected type '*Function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.Name != "f" {
		t.Errorf("expected Function named '%s'; got '%s'", "f", f.Name)
		return
	}
	if l := len(f.Args); l != 0 {
		t.Errorf("expected '%d' argument(s); got '%d'", 2, l)
		return
	}
}

func TestNestedFunctions(t *testing.T) {
	const expectedFuncName = "func"
	str := fmt.Sprintf(`${
							%s(
								"12345",
								g(
									987
								)
							)
						}`, expectedFuncName)
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*Function)
	if !ok {
		t.Errorf("expected type '*Function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.Name != expectedFuncName {
		t.Errorf("expected Function named '%s'; got '%s'", expectedFuncName, f.Name)
		return
	}
	if l := len(f.Args); l != 2 {
		t.Errorf("expected '%d' argument(s); got '%d'", 2, l)
		return
	}
	arg1, ok := f.Args[0].(*StringIdentity)
	if !ok {
		t.Errorf("expected argument of type '*StringIdentity'; got '%s'", reflect.TypeOf(arg1))
		return
	}
	g, ok := f.Args[1].(*Function)
	if !ok {
		t.Errorf("expected argument of type '*Function'; got '%s'", reflect.TypeOf(g))
		return
	}
	if l := len(g.Args); l != 1 {
		t.Errorf("expected '%d' argument(s); got '%d'", 1, l)
		return
	}
	arg2, ok := g.Args[0].(*IntegerIdentity)
	if !ok {
		t.Errorf("expected argument of type '*IntegerIdentity'; got '%s'", reflect.TypeOf(arg2))
		return
	}
}

func TestEvaluateVar(t *testing.T) {
	const expected = "hello"
	str := `${  var  (  "a"    )  }`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	vars := map[string]interface{}{
		"a": expected,
	}
	e, err := token.Evaluate(vars, &http.Request{})
	if err != nil {
		t.Error(err)
		return
	}
	if e != expected {
		t.Errorf("expected value '%s'; got '%v'", expected, e)
		return
	}
}

func TestVarAsFunctionArg(t *testing.T) {
	const expectedFuncName = "link"
	str := fmt.Sprintf(`${%s(var("some_link"))}`, expectedFuncName)
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*Function)
	if !ok {
		t.Errorf("expected type '*Function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.Name != expectedFuncName {
		t.Errorf("expected Function named '%s'; got '%s'", expectedFuncName, f.Name)
		return
	}
	if l := len(f.Args); l != 1 {
		t.Errorf("expected '%d' argument(s); got '%d'", 1, l)
		return
	}
	arg1, ok := f.Args[0].(*Function)
	if !ok {
		t.Errorf("expected argument of type '*Function'; got '%s'", reflect.TypeOf(arg1))
		return
	}
}

func TestComplexIfElseStatement(t *testing.T) {
	str := fmt.Sprintf(`${and(
							if(not(eq("hello", "world")))
								true
							else
								false,
							true)}`)
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(bool)
	if !ok {
		t.Errorf("expected result of type 'bool'; got '%s'", reflect.TypeOf(e))
		return
	}
	if v != true {
		t.Errorf("expected true; got '%t'", v)
		return
	}
}

func TestIfElseStatement(t *testing.T) {
	const expected = "correct"
	str := fmt.Sprintf(`${
							if (contains("Hello, world!", "world"))
								"%s"
							else
								"wrong"
						}`, expected)
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	f, ok := token.(*IfElse)
	if !ok {
		t.Errorf("expected type '*IfElse'; got '%s'", reflect.TypeOf(token))
		return
	}
	e, err := f.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(string)
	if !ok {
		t.Errorf("expected result of type 'string'; got '%s'", reflect.TypeOf(v))
		return
	}
	if v != expected {
		t.Errorf("expected '%s'; got '%s'", expected, v)
		return
	}
}

func TestEvaluateHttpHeader(t *testing.T) {
	const expected = "application/json"
	str := `${http_header("Content-Type")}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	headers := http.Header{
		"Content-Type": []string{expected},
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{Header: headers})
	if err != nil {
		t.Error(err)
		return
	}
	if e != expected {
		t.Errorf("expected value '%s'; got '%v'", expected, e)
		return
	}
}

func TestRequestURL(t *testing.T) {
	expected, _ := url.Parse("https://examp.lecom/foo?bar=buzz")
	str := "${request_url()}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{URL: expected})
	if err != nil {
		t.Error(err)
		return
	}
	if e != expected.String() {
		t.Errorf("expected value '%s'; got '%v'", expected, e)
		return
	}
}

func TestRegexMatch(t *testing.T) {
	str := `${regex_match("Hello, world!",
							"^.*world.*$")}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(bool)
	if !ok {
		t.Errorf("expected type 'bool'; got '%s'", reflect.TypeOf(e))
		return
	}
	if !v {
		t.Errorf("expected value '%t'; got '%t'", true, v)
		return
	}
}

func TestQuery(t *testing.T) {
	str := `${if (ne(request_url_query("bar"), ""))
					request_url_query()
				else
					"wrong"}`
	const query = "bar=buzz"
	u, _ := url.Parse(fmt.Sprintf("https://examp.lecom/foo?%s", query))
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{URL: u})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(string)
	if !ok {
		t.Errorf("expected type 'string'; got '%v'", reflect.TypeOf(e))
		return
	}
	if v != query {
		t.Errorf("expected value '%s'; got '%s'", query, v)
		return
	}
}

func TestHTTPMethod(t *testing.T) {
	const expected = "GET"
	str := fmt.Sprintf(`${if (eq(request_http_method(), "%s"))
					"ok"
				else
					"wrong"}`, expected)
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{Method: expected})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(string)
	if !ok {
		t.Errorf("expected type 'string'; got '%v'", reflect.TypeOf(e))
		return
	}
	if v != "ok" {
		t.Errorf("expected value '%s'; got '%s'", "ok", v)
		return
	}
}

func TestRequestHost(t *testing.T) {
	const expected = "fak.eurl"
	str := fmt.Sprintf(`${
		if (eq(request_http_host(), "%s"))
			"ok"
		else
			"wrong"
	}`, expected)
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{Host: expected})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(string)
	if !ok {
		t.Errorf("expected type 'string'; got '%v'", reflect.TypeOf(e))
		return
	}
	if v != "ok" {
		t.Errorf("expected value '%s'; got '%s'", "ok", v)
		return
	}
}

func TestArrayInTrue(t *testing.T) {
	str := `${
		if (in(["a", "b", "c"], "b"))
			"ok"
		else
			"wrong"
	}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(string)
	if !ok {
		t.Errorf("expected type 'string'; got '%v'", reflect.TypeOf(e))
		return
	}
	if v != "ok" {
		t.Errorf("expected value '%s'; got '%s'", "ok", v)
		return
	}
}

func TestArrayInFalse(t *testing.T) {
	str := `${
		if (not(in([1, 2, 4], 3)))
			"ok"
		else
			"wrong"
	}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.Evaluate(make(map[string]interface{}), &http.Request{})
	if err != nil {
		t.Error(err)
		return
	}
	v, ok := e.(string)
	if !ok {
		t.Errorf("expected type 'string'; got '%v'", reflect.TypeOf(e))
		return
	}
	if v != "ok" {
		t.Errorf("expected value '%s'; got '%s'", "ok", v)
		return
	}
}
