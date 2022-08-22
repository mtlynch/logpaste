package test_sqlite

import (
	"fmt"

	"github.com/mtlynch/logpaste/random"
	"github.com/mtlynch/logpaste/store"
	"github.com/mtlynch/logpaste/store/sqlite"
)

func New() store.Store {
	return sqlite.New(ephemeralDbURI())
}

func ephemeralDbURI() string {
	name := random.String(10)
	return fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
}
