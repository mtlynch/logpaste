package handlers

import (
	"net/http"
	"strings"
)

// ProxyIPHeadersHandler trusts X-Forwarded-For and X-Real-IP from a reverse proxy.
// Enable it only when that proxy overwrites client-supplied forwarding headers
// and clients cannot reach the server directly.
func ProxyIPHeadersHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var address string
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			address, _, _ = strings.Cut(fwd, ",")
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			address = realIP
		}
		if address = strings.TrimSpace(address); address != "" {
			r.RemoteAddr = address
		}
		next.ServeHTTP(w, r)
	})
}
