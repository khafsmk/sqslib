package mqueue

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func checker(t *testing.T) func(err error) {
	return func(err error) {
		if err != nil {
			t.Helper()
			t.Fatal(err)
		}
	}
}

var transformJSON = cmp.FilterValues(func(x, y []byte) bool {
	return json.Valid(x) && json.Valid(y)
}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
	if err := json.Unmarshal(in, &out); err != nil {
		panic(err) // should never occur given previous filter to ensure valid JSON
	}
	return out
}))

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
