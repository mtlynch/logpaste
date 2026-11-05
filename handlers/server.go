package handlers

import (
	"net/http"

	"github.com/mtlynch/logpaste/limit"
	"github.com/mtlynch/logpaste/store"
	"github.com/mtlynch/logpaste/store/sqlite"
)

type Server interface {
	Router() http.Handler
}

func New(sp SiteProperties, perMinuteLimit int, maxCharLimit int64) Server {
	s := defaultServer{
		router:        http.NewServeMux(),
		store:         sqlite.New("data/store.db"),
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
	router        *http.ServeMux
	store         store.Store
	siteProps     SiteProperties
	ipRateLimiter limit.IPRateLimiter
	maxCharLimit  int64
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() http.Handler {
	return s.router
}
