package handlers

import (
	"github.com/gorilla/mux"

	"github.com/mtlynch/logpaste/limit"
	"github.com/mtlynch/logpaste/store"
	"github.com/mtlynch/logpaste/store/sqlite"
)

type Server interface {
	Router() *mux.Router
}

func New(sp SiteProperties, perMinuteLimit int) Server {
	s := defaultServer{
		router:        mux.NewRouter(),
		store:         sqlite.New(),
		siteProps:     sp,
		ipRateLimiter: limit.New(perMinuteLimit),
	}
	s.routes()
	return s
}

type SiteProperties struct {
	Title      string
	Subtitle   string
	FooterHTML string
	DarkMode   bool
	ShowDocs   bool
}

type defaultServer struct {
	router        *mux.Router
	store         store.Store
	siteProps     SiteProperties
	ipRateLimiter limit.IPRateLimiter
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
