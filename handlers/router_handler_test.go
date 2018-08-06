package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/naighes/imposter/cfg"
)

func TestEmptyRuleSet(t *testing.T) {
	config := cfg.Config{}
	r := httptest.NewRecorder()
	routes, err := NewRouterHandler(&config, nil)
	if err != nil {
		t.Errorf("cannot create a new instance of NewRouterHandler")
	}
	routes.ServeHTTP(r, nil)
	const expected = 404
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}

func TestMatchingRule(t *testing.T) {
	const expected = 200
	rsp := cfg.MatchRsp{Body: "hello", StatusCode: "${200}"}
	def := cfg.MatchDef{RuleExpression: `${
		regex_match(request_url_path(), "^/[0-9]+$")
	}`, Response: &rsp}
	defs := []*cfg.MatchDef{&def}
	config := cfg.Config{Defs: defs}
	r := httptest.NewRecorder()
	routes, err := NewRouterHandler(&config, nil)
	if err != nil {
		t.Errorf("cannot create a new instance of NewRouterHandler: %v", err)
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
	rsp := cfg.MatchRsp{Body: "", StatusCode: "${200}"}
	def := cfg.MatchDef{RuleExpression: `${
		and(
			eq(request_http_method(), "GET"),
			contains(request_url_path(), "bbb")
		)
	}`, Response: &rsp}
	defs := []*cfg.MatchDef{&def}
	config := cfg.Config{Defs: defs}
	r := httptest.NewRecorder()
	routes, err := NewRouterHandler(&config, nil)
	if err != nil {
		t.Errorf("cannot create a new instance of NewRouterHandler")
	}
	url, _ := url.Parse("http://fak.eurl/aaa")
	req := http.Request{Method: "GET", URL: url}
	routes.ServeHTTP(r, &req)
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d: %s", expected, r.Code, r.Body.String())
	}
}
