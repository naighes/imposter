package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ResponseMock struct {
	StatusCode int
	HeadersMap map[string]string
	Body       []byte
}

func (r *ResponseMock) MakeResponse() *http.Response {
	status := strconv.Itoa(r.StatusCode) + " " + http.StatusText(r.StatusCode)
	header := http.Header{}
	for name, value := range r.HeadersMap {
		header.Set(name, value)
	}
	contentLength := len(r.Body)
	header.Set("Content-Length", strconv.Itoa(contentLength))
	res := &http.Response{
		Status:           status,
		StatusCode:       r.StatusCode,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           header,
		Body:             ioutil.NopCloser(bytes.NewReader(r.Body)),
		ContentLength:    int64(contentLength),
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		TLS:              nil,
	}
	if r.StatusCode == http.StatusNoContent || r.StatusCode == http.StatusNotModified {
		if res.ContentLength != 0 {
			res.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
			res.ContentLength = 0
		}
		header.Del("Content-Length")
	}
	return res
}
