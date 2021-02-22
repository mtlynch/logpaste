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
	_, err = ctx.Exec(`
CREATE TABLE IF NOT EXISTS entries (
	id TEXT PRIMARY KEY,
	creation_time TEXT,
	contents TEXT
	)`)
	if err != nil {
		log.Fatalln(err)
	}
	return &db{
		ctx: ctx,
	}
}

func (d db) GetEntry(id string) (string, error) {
	stmt, err := d.ctx.Prepare("SELECT contents FROM entries WHERE id=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var contents string
	err = stmt.QueryRow(id).Scan(&contents)
	if err != nil {
		return "", err
	}
	return contents, nil
}

func (d db) InsertEntry(id string, contents string) error {
	stmt, err := d.ctx.Prepare(`
	INSERT INTO entries(
		id,
		creation_time,
		contents)
	values(?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	t := time.Now().Format(time.RFC3339)

	_, err = stmt.Exec(id, t, contents)
	if err != nil {
		return err
	}
	return nil
}
