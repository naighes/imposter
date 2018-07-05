package main

import (
	"fmt"
	"testing"
)

func TestFuncWithOneStringArg(t *testing.T) {
	const expectedName = "link"
	const expectedArg = "http://fak.eurl"
	input := fmt.Sprintf("${%s(\"%s\")}", expectedName, expectedArg)
	r, args, err := ParseFunc(input)
	if err != nil {
		t.Errorf("HandleFunc raised an error: %s", err)
	}
	if r != expectedName {
		t.Errorf("expected function name '%s'; got '%s'", expectedName, r)
	}
	if l := len(args); l != 1 {
		t.Errorf("expected arguments length of %d; got %d", 1, l)
	}
	if a := args[0]; a != expectedArg {
		t.Errorf("expected argument[0] '%s'; got '%s'", expectedArg, a)
	}
}

func TestFuncWithOneStringArgAndTrailingSpaces(t *testing.T) {
	const expectedName = "link"
	const expectedArg = "http://fak.eurl"
	input := fmt.Sprintf("${%s(  \"%s\"  )}", expectedName, expectedArg)
	r, args, err := ParseFunc(input)
	if err != nil {
		t.Errorf("HandleFunc raised an error: %s", err)
	}
	if r != expectedName {
		t.Errorf("expected function name '%s'; got '%s'", expectedName, r)
	}
	if l := len(args); l != 1 {
		t.Errorf("expected arguments length of %d; got %d", 1, l)
	}
	if a := args[0]; a != expectedArg {
		t.Errorf("expected argument[0] '%s'; got '%s'", expectedArg, a)
	}
}

func TestFuncWithOneStringArgAndEscapedDoubleQuotes(t *testing.T) {
	const expectedName = "link"
	const expectedArg = "http\\\"://fak.eurl"
	input := fmt.Sprintf("${%s(  \"%s\"  )}", expectedName, expectedArg)
	r, args, err := ParseFunc(input)
	if err != nil {
		t.Errorf("HandleFunc raised an error: %s", err)
	}
	if r != expectedName {
		t.Errorf("expected function name '%s'; got '%s'", expectedName, r)
	}
	if l := len(args); l != 1 {
		t.Errorf("expected arguments length of %d; got %d", 1, l)
	}
	if a := args[0]; a != expectedArg {
		t.Errorf("expected argument[0] '%s'; got '%s'", expectedArg, a)
	}
}

func TestFuncWithTwoStringArgs(t *testing.T) {
	const expectedName = "link"
	const expectedArg1 = "http://fak.eurl"
	const expectedArg2 = "http://fa.keurl"
	input := fmt.Sprintf("${%s(\"%s\", \"%s\")}", expectedName, expectedArg1, expectedArg2)
	r, args, err := ParseFunc(input)
	if err != nil {
		t.Errorf("HandleFunc raised an error: %s", err)
	}
	if r != expectedName {
		t.Errorf("expected function name '%s'; got '%s'", expectedName, r)
	}
	if l := len(args); l != 2 {
		t.Errorf("expected arguments length of %d; got %d", 1, l)
	}
	if a1, a2 := args[0], args[1]; a1 != expectedArg1 || a2 != expectedArg2 {
		t.Errorf("expected argument[0] '%s'; got '%s'", expectedArg1, a1)
	}
}
