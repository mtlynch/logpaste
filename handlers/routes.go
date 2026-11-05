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
	s.router.HandleFunc("GET /favicon.ico", notFound())
	s.router.Handle("GET /static/", GetStaticFilesHandler())
	s.router.HandleFunc("GET /{id}", s.pasteGet())
	s.router.HandleFunc("OPTIONS /{$}", s.pasteOptions())
	s.router.Handle("PUT /{$}", s.ipRateLimiter.Limit(s.pastePut()))
	s.router.Handle("POST /{$}", s.ipRateLimiter.Limit(s.pastePost()))
	s.router.HandleFunc("GET /{$}", s.serveIndexPage())
}
