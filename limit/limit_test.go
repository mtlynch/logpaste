package limit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimitClientAddress(t *testing.T) {
	tests := []struct {
		explanation          string
		firstAddress         string
		secondAddress        string
		perMinuteLimit       int
		firstStatusExpected  int
		secondStatusExpected int
	}{
		{
			explanation:          "IPv4 with port and bare address share a key",
			firstAddress:         "192.0.2.1:1000",
			secondAddress:        "192.0.2.1",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusTooManyRequests,
		},
		{
			explanation:          "IPv6 with port and bare address share a key",
			firstAddress:         "[2001:db8::1]:1000",
			secondAddress:        "2001:db8::1",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusTooManyRequests,
		},
		{
			explanation:          "same IPv4 on different ports is limited together",
			firstAddress:         "192.0.2.2:1000",
			secondAddress:        "192.0.2.2:2000",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusTooManyRequests,
		},
		{
			explanation:          "canonical IPv6 variants share a key",
			firstAddress:         "[2001:0db8:0:0:0:0:0:1]:1000",
			secondAddress:        "2001:db8::1",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusTooManyRequests,
		},
		{
			explanation:          "mapped IPv4 shares a key with IPv4",
			firstAddress:         "192.0.2.3:1000",
			secondAddress:        "[::ffff:192.0.2.3]:2000",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusTooManyRequests,
		},
		{
			explanation:          "different IPs have separate limits",
			firstAddress:         "192.0.2.4:1000",
			secondAddress:        "192.0.2.5:1000",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusOK,
		},
		{
			explanation:          "malformed addresses return server errors",
			firstAddress:         "not-an-address",
			secondAddress:        "still-not-an-address",
			perMinuteLimit:       1,
			firstStatusExpected:  http.StatusInternalServerError,
			secondStatusExpected: http.StatusInternalServerError,
		},
		{
			explanation:          "disabled limiter permits addresses regardless of format",
			firstAddress:         "not-an-address",
			secondAddress:        "192.0.2.1",
			perMinuteLimit:       RateLimitingDisabled,
			firstStatusExpected:  http.StatusOK,
			secondStatusExpected: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.explanation, func(t *testing.T) {
			limiter := New(tt.perMinuteLimit)
			handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			{
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.RemoteAddr = tt.firstAddress
				recorder := httptest.NewRecorder()
				handler.ServeHTTP(recorder, req)
				if got, want := recorder.Code, tt.firstStatusExpected; got != want {
					t.Fatalf("first request (%q): got status %d, want %d", tt.firstAddress, got, want)
				}
			}

			{
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.RemoteAddr = tt.secondAddress
				recorder := httptest.NewRecorder()
				handler.ServeHTTP(recorder, req)
				if got, want := recorder.Code, tt.secondStatusExpected; got != want {
					t.Fatalf("second request (%q): got status %d, want %d", tt.secondAddress, got, want)
				}
			}
		})
	}
}
