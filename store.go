package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/naighes/imposter/functions"
)

type Store interface {
	Add(string, *functions.HTTPRsp) bool
	Get(string) *functions.HTTPRsp
	ServeWrite(http.ResponseWriter, *http.Request)
	ServeRead(http.ResponseWriter, *http.Request) bool
}

type inMemoryStore struct {
	entries map[string]*functions.HTTPRsp
	lock    sync.RWMutex
}

func newInMemoryStore() Store {
	lock := sync.RWMutex{}
	entries := make(map[string]*functions.HTTPRsp)
	return &inMemoryStore{entries, lock}
}

func (s *inMemoryStore) Add(url string, rsp *functions.HTTPRsp) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.entries[url]
	s.entries[url] = rsp
	return !ok
}

func (s *inMemoryStore) Get(url string) *functions.HTTPRsp {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if e, ok := s.entries[url]; ok {
		return e
	}
	return nil
}

func (s *inMemoryStore) ServeWrite(w http.ResponseWriter, r *http.Request) {
	var body string
	if r.Body == nil {
		body = ""
	} else {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		body = string(b)
	}
	now := time.Now().Format(http.TimeFormat)
	headers := make(http.Header)
	headers.Set("Last-Modified", now)
	headers.Set("Date", now)
	headers.Set("Content-Type", r.Header.Get("Content-Type"))
	created := s.Add(r.URL.String(), &functions.HTTPRsp{Body: body, Headers: headers})
	if created {
		w.WriteHeader(202)
	} else {
		w.WriteHeader(204)
	}
}

func (s *inMemoryStore) ServeRead(w http.ResponseWriter, r *http.Request) bool {
	rsp := s.Get(r.URL.String())
	if rsp == nil {
		return false
	}
	for k := range rsp.Headers {
		w.Header().Set(k, rsp.Headers.Get(k))
	}
	w.WriteHeader(200)
	if r.Method != "HEAD" {
		fmt.Fprintf(w, rsp.Body)
	}
	return true
}
