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

const charactersPerMiB = 1024 * 1024

func New(sp SiteProperties, perMinuteLimit int, maxPasteMiB int64) Server {
	maxCharLimit := maxPasteMiB * charactersPerMiB
	s := defaultServer{
		router:        mux.NewRouter(),
		store:         sqlite.New(),
		siteProps:     sp,
		ipRateLimiter: limit.New(perMinuteLimit),
		maxCharLimit:  maxCharLimit,
	}
	s.routes()
	return s
}

type SiteProperties struct {
	Title      string
	Subtitle   string
	FooterHTML string
	ShowDocs   bool
}

type defaultServer struct {
	router        *mux.Router
	store         store.Store
	siteProps     SiteProperties
	ipRateLimiter limit.IPRateLimiter
	maxCharLimit  int64
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
