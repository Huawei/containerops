package utils

import "net/http"

// GetRequestScheme Returns the scheme of a http request.
func GetRequestScheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme
}
