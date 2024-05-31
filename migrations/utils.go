package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"path/filepath"
	"runtime"
)

func MigrateUp(db *sql.DB) error {
	return goose.Up(db, getCurrentFileDir())
}

func ResetMigration(db *sql.DB) error {
	return goose.Reset(db, getCurrentFileDir())
}

func getCurrentFileDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
