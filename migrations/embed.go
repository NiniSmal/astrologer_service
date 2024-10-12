package migrations

import (
	"embed"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var FS embed.FS

func UpMigrations(pool *pgxpool.Pool) error {
	// setup database
	db := stdlib.OpenDBFromPool(pool)

	goose.SetBaseFS(FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		if errors.Is(err, goose.ErrAlreadyApplied) {
			return nil
		}
		return err
	}
	return nil
}
