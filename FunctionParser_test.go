package main

import (
	"reflect"
	"testing"
)

func TestIdentity(t *testing.T) {
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

func TestNestedFunctions(t *testing.T) {
	str := "${  f  (  \"12345\"  ,    ${  g  (   987   ) }  )  }"
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
