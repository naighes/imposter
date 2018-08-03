package main

import (
	"bytes"
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
	lock    *sync.RWMutex
	rt      RecordType
}

func newInMemoryStore(rt RecordType) Store {
	lock := sync.RWMutex{}
	entries := make(map[string]*functions.HTTPRsp)
	return &inMemoryStore{entries, &lock, rt}
}

func (s *inMemoryStore) getKey(r *http.Request) string {
	u := r.URL
	var b bytes.Buffer
	if s.rt&Scheme == Scheme && u.Scheme != "" {
		b.WriteString(fmt.Sprintf("%s://", u.Scheme))
	}
	if s.rt&Host == Host {
		b.WriteString(u.Host)
	}
	if s.rt&Path == Path {
		b.WriteString(u.Path)
	}
	if s.rt&Query == Query && u.RawQuery != "" {
		b.WriteString(fmt.Sprintf("?%s", u.RawQuery))
	}
	return b.String()
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
		defer r.Body.Close()
		body = string(b)
	}
	now := time.Now().Format(http.TimeFormat)
	headers := make(http.Header)
	headers.Set("Last-Modified", now)
	headers.Set("Date", now)
	headers.Set("Content-Type", r.Header.Get("Content-Type"))
	created := s.Add(s.getKey(r), &functions.HTTPRsp{Body: body, Headers: headers})
	if created {
		w.WriteHeader(202)
	} else {
		w.WriteHeader(204)
	}
}

func (s *inMemoryStore) ServeRead(w http.ResponseWriter, r *http.Request) bool {
	rsp := s.Get(s.getKey(r))
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
