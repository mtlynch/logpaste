package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"

	migrate "codeberg.org/mtlynch/go-evolutionary-migrate"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

func applyMigrations(ctx *sql.DB) {
	migrationsFS, err := fs.Sub(migrationsFs, "migrations")
	if err != nil {
		log.Fatalf("failed to locate migrations directory: %v", err)
	}

	if err := migrate.Run(context.Background(), ctx, migrationsFS); err != nil {
		log.Fatalf("failed to apply DB migrations: %v", err)
	}
}
