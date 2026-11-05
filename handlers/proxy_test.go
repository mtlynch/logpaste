package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mtlynch/logpaste/limit"
)

func TestProxyIPHeadersHandler(t *testing.T) {
	for _, tt := range []struct {
		explanation     string
		forwardedFor    string
		realIP          string
		forwarded       string
		addressExpected string
	}{
		{
			explanation:     "requests without proxy headers retain the peer address",
			addressExpected: "192.0.2.1:1234",
		},
		{
			explanation:     "forwarded-for uses the first address and takes precedence",
			forwardedFor:    " 198.51.100.1 ,203.0.113.2",
			realIP:          "198.51.100.2",
			forwarded:       "for=198.51.100.3",
			addressExpected: "198.51.100.1",
		},
		{
			explanation:     "real-IP supplies the client address when forwarded-for is absent",
			realIP:          " 198.51.100.2 ",
			forwarded:       "for=198.51.100.3",
			addressExpected: "198.51.100.2",
		},
		{
			explanation:     "Forwarded alone is ignored",
			forwarded:       "for=198.51.100.3;proto=https",
			addressExpected: "192.0.2.1:1234",
		},
	} {
		t.Run(tt.explanation, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.RemoteAddr = "192.0.2.1:1234"
			r.Header.Set("X-Forwarded-For", tt.forwardedFor)
			r.Header.Set("X-Real-IP", tt.realIP)
			r.Header.Set("Forwarded", tt.forwarded)
			r.Header.Set("X-Forwarded-Host", "untrusted.example")
			r.Header.Set("X-Forwarded-Proto", "https")
			h := ProxyIPHeadersHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got, want := r.RemoteAddr, tt.addressExpected; got != want {
					t.Errorf("RemoteAddr=%q, want=%q", got, want)
				}
				if got, want := r.Host, "example.com"; got != want {
					t.Errorf("Host=%q, want=%q", got, want)
				}
				if got, want := r.URL.Scheme, ""; got != want {
					t.Errorf("scheme=%q, want=%q", got, want)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if got, want := w.Code, http.StatusNoContent; got != want {
				t.Fatalf("status=%d, want=%d", got, want)
			}
		})
	}
}

func TestProxyIPRateLimiting(t *testing.T) {
	limiter := limit.New(1)
	h := ProxyIPHeadersHandler(limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})))
	{
		r := httptest.NewRequest(http.MethodPut, "/", nil)
		r.Header.Set("X-Forwarded-For", "2001:db8::1, 192.0.2.1")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if got, want := w.Code, http.StatusNoContent; got != want {
			t.Fatalf("status=%d, want=%d", got, want)
		}
	}
	{
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set("X-Real-IP", "2001:db8::1")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if got, want := w.Code, http.StatusTooManyRequests; got != want {
			t.Fatalf("status=%d, want=%d", got, want)
		}
	}
	{
		r := httptest.NewRequest(http.MethodPut, "/", nil)
		r.Header.Set("X-Forwarded-For", "2001:db8::2")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if got, want := w.Code, http.StatusNoContent; got != want {
			t.Fatalf("status=%d, want=%d", got, want)
		}
	}
}
