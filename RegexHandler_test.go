package main

import (
	"net/http/httptest"
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
