package mqueue

import (
	"io"
	"net/http"
	"strings"
)

type clientDoFunc func(*http.Request) (*http.Response, error)

func (f clientDoFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

var nopClient = clientDoFunc(func(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
})
