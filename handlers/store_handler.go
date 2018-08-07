package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/naighes/imposter/functions"
)

type recordType int

const (
	scheme recordType = 1 << iota
	host
	path
	query
)

// StoreHandler type defines an handler to support imPOSTer recording capabilities.
type StoreHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) bool
}

type inMemoryStoreHandler struct {
	entries map[string]*functions.HTTPRsp
	lock    *sync.RWMutex
	rt      recordType
}

func parseRecordConfig(config string) (recordType, error) {
	if config == "" {
		return 0, nil
	}
	m := map[string]recordType{
		"scheme": scheme,
		"host":   host,
		"path":   path,
		"query":  query,
	}
	var r recordType
	e := strings.Split(config, "|")
	for _, v := range e {
		rt, ok := m[v]
		if !ok {
			return 0, fmt.Errorf("'%s' is not a valid flag: select multiple values from {'scheme', 'host', 'path', 'query'} separated by pipe (|)", v)
		}
		r = r | rt
	}
	return r, nil
}

// NewInMemoryStoreHandler builds a new instance of StoreHandler.
func NewInMemoryStoreHandler(config string) (StoreHandler, error) {
	lock := sync.RWMutex{}
	entries := make(map[string]*functions.HTTPRsp)
	rt, err := parseRecordConfig(config)
	if err != nil {
		return nil, err
	}
	return &inMemoryStoreHandler{entries, &lock, rt}, nil
}

func (s *inMemoryStoreHandler) getKey(r *http.Request) string {
	u := r.URL
	var b bytes.Buffer
	if s.rt&scheme == scheme && u.Scheme != "" {
		b.WriteString(fmt.Sprintf("%s://", u.Scheme))
	}
	if s.rt&host == host {
		b.WriteString(u.Host)
	}
	if s.rt&path == path {
		b.WriteString(u.Path)
	}
	if s.rt&query == query && u.RawQuery != "" {
		b.WriteString(fmt.Sprintf("?%s", u.RawQuery))
	}
	return b.String()
}

func (s *inMemoryStoreHandler) add(url string, rsp *functions.HTTPRsp) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.entries[url]
	s.entries[url] = rsp
	return !ok
}

func (s *inMemoryStoreHandler) get(url string) *functions.HTTPRsp {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if e, ok := s.entries[url]; ok {
		return e
	}
	return nil
}

func (s *inMemoryStoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	switch r.Method {
	case "PUT":
		return s.serveWrite(w, r)
	case "GET", "HEAD":
		return s.serveRead(w, r)
	default:
		return false
	}
}

func (s *inMemoryStoreHandler) serveWrite(w http.ResponseWriter, r *http.Request) bool {
	var body string
	if r.Body == nil {
		body = ""
	} else {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return false
		}
		defer r.Body.Close()
		body = string(b)
	}
	now := time.Now().Format(http.TimeFormat)
	headers := make(http.Header)
	headers.Set("Last-Modified", now)
	headers.Set("Date", now)
	headers.Set("Content-Type", r.Header.Get("Content-Type"))
	created := s.add(s.getKey(r), &functions.HTTPRsp{Body: body, Headers: headers})
	if created {
		w.WriteHeader(202)
	} else {
		w.WriteHeader(204)
	}
	return true
}

func (s *inMemoryStoreHandler) serveRead(w http.ResponseWriter, r *http.Request) bool {
	rsp := s.get(s.getKey(r))
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
