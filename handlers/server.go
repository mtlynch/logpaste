package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mtlynch/logpaste/store"
	"github.com/mtlynch/logpaste/store/sqlite"
)

type Server interface {
	Router() *mux.Router
}

func New(sp SiteProperties) Server {
	s := defaultServer{
		router:    mux.NewRouter(),
		store:     sqlite.New(),
		siteProps: sp,
	}
	s.routes()
	return s
}

type httpMiddlewareHandler func(http.Handler) http.Handler

type SiteProperties struct {
	Title    string
	Subtitle string
}

type defaultServer struct {
	router    *mux.Router
	store     store.Store
	siteProps SiteProperties
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
