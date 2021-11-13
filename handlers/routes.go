package handlers

import (
	"net/http"
)

func notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Resource not found", http.StatusNotFound)
	}
}

func (s *defaultServer) routes() {
	s.router.HandleFunc("/favicon.ico", notFound()).Methods(http.MethodGet)
	s.router.HandleFunc("/css/", s.serveStaticResource()).Methods(http.MethodGet)
	s.router.HandleFunc("/js/", s.serveStaticResource()).Methods(http.MethodGet)
	s.router.HandleFunc("/third-party/", s.serveStaticResource()).Methods(http.MethodGet)
	s.router.HandleFunc("/{id}", s.pasteGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/", s.pasteOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/", s.ipRateLimiter.Limit(s.pastePut())).Methods(http.MethodPut)
	s.router.HandleFunc("/", s.ipRateLimiter.Limit(s.pastePost())).Methods(http.MethodPost)
	s.router.HandleFunc("/", s.serveIndexPage()).Methods(http.MethodGet)
}
