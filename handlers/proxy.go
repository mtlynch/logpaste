package handlers

import (
	"net/http"
	"regexp"
	"strings"
)

var forwardedForPattern = regexp.MustCompile(`(?i)(?:^|[;,]\s*)for=(?:"([^"]*)"|([^;,\s]+))`)

// ProxyIPHeadersHandler trusts client IP headers supplied by a reverse proxy.
// Enable it only when that proxy overwrites client-supplied forwarding headers
// and clients cannot reach the server directly.
func ProxyIPHeadersHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var address string
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			address, _, _ = strings.Cut(fwd, ",")
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			address = realIP
		} else if match := forwardedForPattern.FindStringSubmatch(r.Header.Get("Forwarded")); match != nil {
			address = match[1]
			if address == "" {
				address = match[2]
			}
			// RFC Forwarded brackets IPv6 literals, even without a port.
			if strings.HasPrefix(address, "[") && strings.HasSuffix(address, "]") {
				address = address[1 : len(address)-1]
			}
		}
		if address = strings.TrimSpace(address); address != "" {
			r.RemoteAddr = address
		}
		next.ServeHTTP(w, r)
	})
}
