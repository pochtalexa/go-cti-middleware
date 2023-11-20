package migrations

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var SQLFiles embed.FS

func ApplyMigrations(DBconn *sql.DB) error {
	fsys := SQLFiles

	goose.SetBaseFS(fsys)
	goose.SetSequential(true)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(DBconn, "."); err != nil {
		return err
	}

	return nil
}
