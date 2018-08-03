package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestReadStoredURL(t *testing.T) {
	store := newInMemoryStore(Scheme | Host | Path | Query)
	u, _ := url.Parse("http://fak.eurl")
	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	store.ServeWrite(r1, &http.Request{URL: u})
	exists := store.ServeRead(r2, &http.Request{URL: u})
	const expected = 200
	if r2.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r2.Code)
	}
	if !exists {
		t.Errorf("an esisting resorce was expected")
	}
}

func TestReadNonStoredURL(t *testing.T) {
	store := newInMemoryStore(Scheme | Host | Path | Query)
	u, _ := url.Parse("http://fak.eurl")
	r := httptest.NewRecorder()
	exists := store.ServeRead(r, &http.Request{URL: u})
	if exists {
		t.Errorf("an esisting resorce was not expected")
	}
}

func TestWriteStoredURL(t *testing.T) {
	store := newInMemoryStore(Scheme | Host | Path | Query)
	u, _ := url.Parse("http://fak.eurl")
	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	store.ServeWrite(r1, &http.Request{URL: u})
	store.ServeWrite(r2, &http.Request{URL: u})
	const expected = 204
	if r2.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r2.Code)
	}
}

func TestWriteNonStoredURL(t *testing.T) {
	store := newInMemoryStore(Scheme | Host | Path | Query)
	u, _ := url.Parse("http://fak.eurl")
	r := httptest.NewRecorder()
	store.ServeWrite(r, &http.Request{URL: u})
	const expected = 202
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}

func TestReadStoredURLWithoutQueryFlag(t *testing.T) {
	store := newInMemoryStore(Scheme | Host | Path)
	u1, _ := url.Parse("http://fak.eurl?id=1")
	u2, _ := url.Parse("http://fak.eurl?id=2")
	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	store.ServeWrite(r1, &http.Request{URL: u1})
	exists := store.ServeRead(r2, &http.Request{URL: u2})
	const expected = 200
	if r2.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r2.Code)
	}
	if !exists {
		t.Errorf("an esisting resorce was expected")
	}
}
