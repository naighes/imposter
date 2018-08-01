package main

import (
	"testing"

	"github.com/naighes/imposter/functions"
)

func TestNonBooleanRuleExpression(t *testing.T) {
	rsp := MatchRsp{}
	def := &MatchDef{RuleExpression: `${"some string"}`, Response: &rsp}
	vars := make(map[string]interface{})
	errors := def.validate(functions.ParseExpression, vars)
	const expected = 1
	if l := len(errors); l != expected {
		t.Errorf("expected %d error(s); got %d instead", expected, l)
		return
	}
}

func TestBooleanRuleExpression(t *testing.T) {
	rsp := MatchRsp{}
	def := &MatchDef{RuleExpression: `${true}`, Response: &rsp}
	vars := make(map[string]interface{})
	errors := def.validate(functions.ParseExpression, vars)
	const expected = 0
	if l := len(errors); l != expected {
		t.Errorf("expected %d error(s); got %d instead", expected, l)
		return
	}
}

func TestBodyExpressionSyntaxError(t *testing.T) {
	rsp := MatchRsp{Body: `${var("www"}`}
	def := &MatchDef{RuleExpression: `${true}`, Response: &rsp}
	vars := make(map[string]interface{})
	errors := def.validate(functions.ParseExpression, vars)
	const expected = 1
	if l := len(errors); l != expected {
		t.Errorf("expected %d error(s); got %d instead", expected, l)
		return
	}
}

func TestHeaderExpressionSyntaxError(t *testing.T) {
	h := map[string]interface{}{
		"Content-Type": "${www}",
	}
	rsp := MatchRsp{Headers: h}
	def := &MatchDef{RuleExpression: `${true}`, Response: &rsp}
	vars := make(map[string]interface{})
	errors := def.validate(functions.ParseExpression, vars)
	const expected = 1
	if l := len(errors); l != expected {
		t.Errorf("expected %d error(s); got %d instead", expected, l)
		return
	}
}
