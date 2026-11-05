package limit

import (
	"log"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const RateLimitingDisabled = 0

type IPRateLimiter struct {
	limiters       map[string]*rate.Limiter
	mu             *sync.Mutex
	perMinuteLimit int
}

func New(perMinuteLimit int) IPRateLimiter {
	iprl := IPRateLimiter{
		limiters:       make(map[string]*rate.Limiter),
		mu:             &sync.Mutex{},
		perMinuteLimit: perMinuteLimit,
	}
	return iprl
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise create a new rate limiter and add it to
// the visitors map, using the IP address as the key.
func (iprl IPRateLimiter) getRateLimiter(ip string) *rate.Limiter {
	iprl.mu.Lock()
	defer iprl.mu.Unlock()

	limiter, exists := iprl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute), iprl.perMinuteLimit)
		iprl.limiters[ip] = limiter
	}

	return limiter
}

func (iprl IPRateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if iprl.perMinuteLimit == RateLimitingDisabled {
			next.ServeHTTP(w, r)
			return
		}

		// Get the IP address for the current user.
		addrPort, err := netip.ParseAddrPort(r.RemoteAddr)
		var addr netip.Addr
		if err == nil {
			addr = addrPort.Addr()
		} else {
			addr, err = netip.ParseAddr(r.RemoteAddr)
		}
		if err != nil {
			log.Printf("Error retrieving user IP: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		limiter := iprl.getRateLimiter(addr.Unmap().String())
		if !limiter.Allow() {
			http.Error(w, "Too many requests. Wait a minute and try again", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
