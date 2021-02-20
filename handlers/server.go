package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server interface {
	Router() *mux.Router
}

func New() Server {
	s := defaultServer{
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

type httpMiddlewareHandler func(http.Handler) http.Handler

type defaultServer struct {
	router *mux.Router
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
