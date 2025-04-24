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

func New(perMinuteLimit int, maxCharLimit int64) Server {
	s := defaultServer{
		router:        mux.NewRouter(),
		store:         sqlite.New(),
		ipRateLimiter: limit.New(perMinuteLimit),
		maxCharLimit:  maxCharLimit,
	}
	s.routes()
	return s
}

type defaultServer struct {
	router        *mux.Router
	store         store.Store
	ipRateLimiter limit.IPRateLimiter
	maxCharLimit  int64
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
