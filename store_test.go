package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/naighes/imposter/functions"
)

func TestReadStoredURL(t *testing.T) {
	store := newInMemoryStore()
	u, _ := url.Parse("http://fak.eurl")
	rsp := functions.HTTPRsp{}
	store.Add(u.String(), &rsp)
	r := httptest.NewRecorder()
	exists := store.ServeRead(r, &http.Request{URL: u})
	const expected = 200
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
	if !exists {
		t.Errorf("an esisting resorce was expected")
	}
}

func TestReadNonStoredURL(t *testing.T) {
	store := newInMemoryStore()
	u, _ := url.Parse("http://fak.eurl")
	r := httptest.NewRecorder()
	exists := store.ServeRead(r, &http.Request{URL: u})
	if exists {
		t.Errorf("an esisting resorce was not expected")
	}
}

func TestWriteStoredURL(t *testing.T) {
	store := newInMemoryStore()
	u, _ := url.Parse("http://fak.eurl")
	rsp := functions.HTTPRsp{}
	store.Add(u.String(), &rsp)
	r := httptest.NewRecorder()
	store.ServeWrite(r, &http.Request{URL: u})
	const expected = 204
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}

func TestWriteNonStoredURL(t *testing.T) {
	store := newInMemoryStore()
	u, _ := url.Parse("http://fak.eurl")
	r := httptest.NewRecorder()
	store.ServeWrite(r, &http.Request{URL: u})
	const expected = 202
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}
