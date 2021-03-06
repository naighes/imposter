package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestReadStoredURL(t *testing.T) {
	store, err := NewInMemoryStoreHandler("scheme|host|path|query")
	if err != nil {
		t.Errorf("cannot create a new instance of InMemoryStoreHandler: %v", err)
	}
	u, _ := url.Parse("http://fak.eurl")
	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	store.ServeHTTP(r1, &http.Request{URL: u, Method: "PUT"})
	exists := store.ServeHTTP(r2, &http.Request{URL: u, Method: "GET"})
	const expected = 200
	if r2.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r2.Code)
	}
	if !exists {
		t.Errorf("an esisting resorce was expected")
	}
}

func TestReadNonStoredURL(t *testing.T) {
	store, err := NewInMemoryStoreHandler("scheme|host|path|query")
	if err != nil {
		t.Errorf("cannot create a new instance of InMemoryStoreHandler: %v", err)
	}
	u, _ := url.Parse("http://fak.eurl")
	r := httptest.NewRecorder()
	exists := store.ServeHTTP(r, &http.Request{URL: u, Method: "GET"})
	if exists {
		t.Errorf("an esisting resorce was not expected")
	}
}

func TestWriteStoredURL(t *testing.T) {
	store, err := NewInMemoryStoreHandler("scheme|host|path|query")
	if err != nil {
		t.Errorf("cannot create a new instance of InMemoryStoreHandler: %v", err)
	}
	u, _ := url.Parse("http://fak.eurl")
	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	store.ServeHTTP(r1, &http.Request{URL: u, Method: "PUT"})
	store.ServeHTTP(r2, &http.Request{URL: u, Method: "PUT"})
	const expected = 204
	if r2.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r2.Code)
	}
}

func TestWriteNonStoredURL(t *testing.T) {
	store, err := NewInMemoryStoreHandler("scheme|host|path|query")
	if err != nil {
		t.Errorf("cannot create a new instance of InMemoryStoreHandler: %v", err)
	}
	u, _ := url.Parse("http://fak.eurl")
	r := httptest.NewRecorder()
	store.ServeHTTP(r, &http.Request{URL: u, Method: "PUT"})
	const expected = 202
	if r.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r.Code)
	}
}

func TestReadStoredURLWithoutQueryFlag(t *testing.T) {
	store, err := NewInMemoryStoreHandler("scheme|host|path")
	if err != nil {
		t.Errorf("cannot create a new instance of InMemoryStoreHandler: %v", err)
	}
	u1, _ := url.Parse("http://fak.eurl?id=1")
	u2, _ := url.Parse("http://fak.eurl?id=2")
	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	store.ServeHTTP(r1, &http.Request{URL: u1, Method: "PUT"})
	exists := store.ServeHTTP(r2, &http.Request{URL: u2, Method: "GET"})
	const expected = 200
	if r2.Code != expected {
		t.Errorf("expected status code %d; got %d", expected, r2.Code)
	}
	if !exists {
		t.Errorf("an esisting resorce was expected")
	}
}
