package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEmptyRuleSet(t *testing.T) {
	config := Config{}
	r := httptest.NewRecorder()
	routes, err := NewRegexHandler(&config)
	if err != nil {
		t.Errorf("cannot create a new instance of NewRegexHandler")
	}
	routes.ServeHTTP(r, nil)
	const expected = 404
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}

func TestMatchingRule(t *testing.T) {
	const expected = 200
	rsp := MatchRsp{Body: "hello", StatusCode: expected}
	def := MatchDef{RuleExpression: `${
		regex_match(request_url_path(), "^/[0-9]+$")
	}`, Response: &rsp}
	defs := []*MatchDef{&def}
	config := Config{Defs: defs}
	r := httptest.NewRecorder()
	routes, err := NewRegexHandler(&config)
	if err != nil {
		t.Errorf("cannot create a new instance of NewRegexHandler: %v", err)
	}
	url, _ := url.Parse("http://fak.eurl/123")
	req := http.Request{Method: "GET", URL: url}
	routes.ServeHTTP(r, &req)
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}

func TestNonMatchingRule(t *testing.T) {
	const expected = 404
	rsp := MatchRsp{Body: "", StatusCode: 200}
	def := MatchDef{RuleExpression: `${
		and(
			eq(request_http_method(), "GET"),
			contains(request_url_path(), "bbb")
		)
	}`, Response: &rsp}
	defs := []*MatchDef{&def}
	config := Config{Defs: defs}
	r := httptest.NewRecorder()
	routes, err := NewRegexHandler(&config)
	if err != nil {
		t.Errorf("cannot create a new instance of NewRegexHandler")
	}
	url, _ := url.Parse("http://fak.eurl/aaa")
	req := http.Request{Method: "GET", URL: url}
	routes.ServeHTTP(r, &req)
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}