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
	ensureDirExists(dbDir)
	ctx, err := sql.Open("sqlite3", dbDir+"/store.db")
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := ctx.Exec(`
-- Apply Litestream recommendations: https://litestream.io/tips/
PRAGMA busy_timeout = 5000;
PRAGMA synchronous = NORMAL;
PRAGMA journal_mode = WAL;
PRAGMA wal_autocheckpoint = 0;
		`); err != nil {
		log.Fatalf("failed to set pragmas: %v", err)
	}

	applyMigrations(ctx)

	return &db{
		ctx: ctx,
	}
}

func (d db) GetEntry(id string) (string, error) {
	var contents string
	if err := d.ctx.QueryRow("SELECT contents FROM entries WHERE id=?", id).Scan(&contents); err != nil {
		if err == sql.ErrNoRows {
			return "", store.EntryNotFoundError{ID: id}
		}
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

func ensureDirExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
