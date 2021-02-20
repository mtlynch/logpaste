package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	s.router.PathPrefix("/{id}").HandlerFunc(s.pasteGet()).Methods(http.MethodGet)
	s.router.PathPrefix("/").HandlerFunc(s.pasteOptions()).Methods(http.MethodOptions)
	s.router.PathPrefix("/").HandlerFunc(s.pastePut()).Methods(http.MethodPut)
	s.router.PathPrefix("/").HandlerFunc(s.serveIndexPage()).Methods(http.MethodGet)
}
