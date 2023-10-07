package handlers

import (
	"net/http"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/images/favicon.ico")
}

func (s *defaultServer) routes() {
	s.router.HandleFunc("/favicon.ico", faviconHandler).Methods(http.MethodGet)
	s.router.PathPrefix("/css/").HandlerFunc(serveStaticResource()).Methods(http.MethodGet)
	s.router.PathPrefix("/js/").HandlerFunc(serveStaticResource()).Methods(http.MethodGet)
	s.router.PathPrefix("/third-party/").HandlerFunc(serveStaticResource()).Methods(http.MethodGet)
	s.router.PathPrefix("/{id}").HandlerFunc(s.pasteGet()).Methods(http.MethodGet)
	s.router.PathPrefix("/").HandlerFunc(s.pasteOptions()).Methods(http.MethodOptions)
	s.router.PathPrefix("/").Handler(s.ipRateLimiter.Limit(s.pastePut())).Methods(http.MethodPut)
	s.router.PathPrefix("/").Handler(s.ipRateLimiter.Limit(s.pastePost())).Methods(http.MethodPost)
	s.router.PathPrefix("/").HandlerFunc(s.serveIndexPage()).Methods(http.MethodGet)
}
