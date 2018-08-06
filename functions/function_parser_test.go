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
	_, ok := token.(*integerIdentity)
	if !ok {
		t.Errorf("expected type '*integerIdentity'; got '%s'", reflect.TypeOf(token))
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
	_, ok := token.(*floatIdentity)
	if !ok {
		t.Errorf("expected type '*floatIdentity'; got '%s'", reflect.TypeOf(token))
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
	_, ok := token.(*stringIdentity)
	if !ok {
		t.Errorf("expected type '*stringIdentity'; got '%s'", reflect.TypeOf(token))
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
	b, ok := token.(*boolIdentity)
	if !ok {
		t.Errorf("expected type '*boolIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := b.Evaluate(ctx)
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
	f, ok := token.(*arrayIdentity)
	if !ok {
		t.Errorf("expected type '*arrayIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := f.Evaluate(ctx)
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
	f, ok := token.(*arrayIdentity)
	if !ok {
		t.Errorf("expected type '*arrayIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	_, err = f.Evaluate(ctx)
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
	b, ok := token.(*function)
	if !ok {
		t.Errorf("expected type '*function'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := b.Evaluate(ctx)
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
	b, ok := token.(*function)
	if !ok {
		t.Errorf("expected type '*function'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := b.Evaluate(ctx)
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
	_, ok := token.(*stringIdentity)
	if !ok {
		t.Errorf("expected type '*stringIdentity'; got '%s'", reflect.TypeOf(token))
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
	si, ok := token.(*stringIdentity)
	if !ok {
		t.Errorf("expected type '*stringIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	s, err := si.Evaluate(ctx)
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
	f, ok := token.(*function)
	if !ok {
		t.Errorf("expected type '*function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.name != "f" {
		t.Errorf("expected Function named '%s'; got '%s'", "f", f.name)
		return
	}
	if l := len(f.args); l != 2 {
		t.Errorf("expected '%d' argument(s); got '%d'", 2, l)
		return
	}
	arg1, ok := f.args[0].(*stringIdentity)
	if !ok {
		t.Errorf("expected argument of type '*stringIdentity'; got '%s'", reflect.TypeOf(arg1))
		return
	}
	arg2, ok := f.args[1].(*integerIdentity)
	if !ok {
		t.Errorf("expected argument of type '*integerIdentity'; got '%s'", reflect.TypeOf(arg2))
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
	f, ok := token.(*function)
	if !ok {
		t.Errorf("expected type '*function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.name != "f" {
		t.Errorf("expected Function named '%s'; got '%s'", "f", f.name)
		return
	}
	if l := len(f.args); l != 0 {
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
	f, ok := token.(*function)
	if !ok {
		t.Errorf("expected type '*function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.name != expectedFuncName {
		t.Errorf("expected Function named '%s'; got '%s'", expectedFuncName, f.name)
		return
	}
	if l := len(f.args); l != 2 {
		t.Errorf("expected '%d' argument(s); got '%d'", 2, l)
		return
	}
	arg1, ok := f.args[0].(*stringIdentity)
	if !ok {
		t.Errorf("expected argument of type '*stringIdentity'; got '%s'", reflect.TypeOf(arg1))
		return
	}
	g, ok := f.args[1].(*function)
	if !ok {
		t.Errorf("expected argument of type '*function'; got '%s'", reflect.TypeOf(g))
		return
	}
	if l := len(g.args); l != 1 {
		t.Errorf("expected '%d' argument(s); got '%d'", 1, l)
		return
	}
	arg2, ok := g.args[0].(*integerIdentity)
	if !ok {
		t.Errorf("expected argument of type '*integerIdentity'; got '%s'", reflect.TypeOf(arg2))
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
	ctx := &EvaluationContext{Vars: vars, Req: &http.Request{}}
	e, err := token.Evaluate(ctx)
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
	f, ok := token.(*function)
	if !ok {
		t.Errorf("expected type '*function'; got '%s'", reflect.TypeOf(token))
		return
	}
	if f.name != expectedFuncName {
		t.Errorf("expected Function named '%s'; got '%s'", expectedFuncName, f.name)
		return
	}
	if l := len(f.args); l != 1 {
		t.Errorf("expected '%d' argument(s); got '%d'", 1, l)
		return
	}
	arg1, ok := f.args[0].(*function)
	if !ok {
		t.Errorf("expected argument of type '*function'; got '%s'", reflect.TypeOf(arg1))
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := token.Evaluate(ctx)
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
	f, ok := token.(*ifElse)
	if !ok {
		t.Errorf("expected type '*ifElse'; got '%s'", reflect.TypeOf(token))
		return
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := f.Evaluate(ctx)
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

func TestEvaluateHTTPHeader(t *testing.T) {
	const expected = "application/json"
	str := `${request_http_header("Content-Type")}`
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	headers := http.Header{
		"Content-Type": []string{expected},
	}
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{Header: headers}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{URL: expected}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{URL: u}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{Method: expected}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{Host: expected}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := token.Evaluate(ctx)
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
	ctx := &EvaluationContext{Vars: make(map[string]interface{}), Req: &http.Request{}}
	e, err := token.Evaluate(ctx)
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
