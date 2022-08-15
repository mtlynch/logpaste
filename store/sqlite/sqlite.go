package sqlite

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mtlynch/logpaste/store"
)

type db struct {
	ctx *sql.DB
}

func New() store.Store {
	dbDir := "data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}
	ctx, err := sql.Open("sqlite3", dbDir+"/store.db")
	if err != nil {
		log.Fatalln(err)
	}

	initStmts := []string{
		// The Litestream documentation recommends these pragmas.
		// https://litestream.io/tips/
		`PRAGMA busy_timeout = 5000`,
		`PRAGMA synchronous = NORMAL`,

		`CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			creation_time TEXT,
			contents TEXT
			)`,
	}
	for _, stmt := range initStmts {
		_, err = ctx.Exec(stmt)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &db{
		ctx: ctx,
	}
}

func (d db) GetEntry(id string) (string, error) {
	var contents string
	if err := d.ctx.QueryRow("SELECT contents FROM entries WHERE id=?", id).Scan(&contents); err != nil {
		return "", err
	}
	return contents, nil
}

func (d db) InsertEntry(id string, contents string) error {
	_, err := d.ctx.Exec(`
	INSERT INTO entries(
		id,
		creation_time,
		contents)
	values(?,?,?)`, id, time.Now().Format(time.RFC3339), contents)
	return err
}
