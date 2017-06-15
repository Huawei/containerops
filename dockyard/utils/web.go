package utils

import "net/http"

// GetRequestScheme Returns the scheme of a http request.
func GetRequestScheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	} else { // Request from a proxy server to unix socket, read the proto from header
		scheme = r.Header.Get("X-Forwarded-Proto")
	}
	return scheme
}
