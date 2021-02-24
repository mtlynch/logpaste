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
	s.router.PathPrefix("/css/").HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)
	s.router.PathPrefix("/js/").HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)
	s.router.PathPrefix("/{id}").HandlerFunc(s.pasteGet()).Methods(http.MethodGet)
	s.router.PathPrefix("/").HandlerFunc(s.pasteOptions()).Methods(http.MethodOptions)
	s.router.PathPrefix("/").HandlerFunc(s.pastePut()).Methods(http.MethodPut)
	s.router.PathPrefix("/").HandlerFunc(s.serveIndexPage()).Methods(http.MethodGet)
}
