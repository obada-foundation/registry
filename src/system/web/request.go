package web

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

// Query returns the argument value from the query string.
func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(val)
}
