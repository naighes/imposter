package main

import (
	"fmt"
	"net/http"
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

func TestBooleanIdentity(t *testing.T) {
	str := "${true}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	b, ok := token.(*booleanIdentity)
	if !ok {
		t.Errorf("expected type '*booleanIdentity'; got '%s'", reflect.TypeOf(token))
		return
	}
	e, err := b.evaluate(make(map[string]interface{}), &http.Request{})
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
		t.Errorf("expected boolean value 'true'; got '%v' instead", v)
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
	e, err := b.evaluate(make(map[string]interface{}), &http.Request{})
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
		t.Errorf("expected boolean value 'true'; got '%v' instead", v)
		return
	}
}

func TestAndEvaluation(t *testing.T) {
	str := "${and(true, true)}"
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
	e, err := b.evaluate(make(map[string]interface{}), &http.Request{})
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
		t.Errorf("expected boolean value 'true'; got '%v' instead", v)
		return
	}
}

func TestBlockWithJustIdentity(t *testing.T) {
	str := "${\"abc\"}"
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
	s, err := si.evaluate(make(map[string]interface{}), &http.Request{})
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
	str := "${  f  (  \"12345\"  ,    987   )  }"
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
		t.Errorf("expected function named '%s'; got '%s'", "f", f.name)
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
		t.Errorf("expected function named '%s'; got '%s'", "f", f.name)
		return
	}
	if l := len(f.args); l != 0 {
		t.Errorf("expected '%d' argument(s); got '%d'", 2, l)
		return
	}
}

func TestNestedFunctions(t *testing.T) {
	const expectedFuncName = "func"
	str := fmt.Sprintf("${  %s  (  \"12345\"  ,      g  (   987   )   )  }", expectedFuncName)
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
		t.Errorf("expected function named '%s'; got '%s'", expectedFuncName, f.name)
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
	str := "${  var  (  \"a\"    )  }"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	vars := map[string]interface{}{
		"a": expected,
	}
	e, err := token.evaluate(vars, &http.Request{})
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
	str := fmt.Sprintf("${%s(var(\"some_link\"))}", expectedFuncName)
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
		t.Errorf("expected function named '%s'; got '%s'", expectedFuncName, f.name)
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
	str := fmt.Sprintf("${ and  (  if (  ne(\"hello\", \"world\")   )    true    else     false  , true) }")
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	e, err := token.evaluate(make(map[string]interface{}), &http.Request{})
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
	str := fmt.Sprintf("${   if (  contains(\"Hello, world!\", \"world\")   )    \"%s\"    else     \"wrong\"   }", expected)
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
	e, err := f.evaluate(make(map[string]interface{}), &http.Request{})
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
	str := "${http_header(\"Content-Type\")}"
	token, err := ParseExpression(str)
	if err != nil {
		t.Error(err)
		return
	}
	headers := http.Header{
		"Content-Type": []string{expected},
	}
	e, err := token.evaluate(make(map[string]interface{}), &http.Request{Header: headers})
	if err != nil {
		t.Error(err)
		return
	}
	if e != expected {
		t.Errorf("expected value '%s'; got '%v'", expected, e)
		return
	}
}
