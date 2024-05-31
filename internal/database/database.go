package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	Driver     = "postgres"
	DBUser     = "test_user"
	DBPassword = "secret"
	Host       = "localhost"
	DBPort     = 52345
)

func ConnectDB() (*sql.DB, error) {
	DSN := fmt.Sprintf(
		"%s://%s:%s@%s:%d/default?sslmode=disable",
		Driver,
		DBUser,
		DBPassword,
		Host,
		DBPort,
	)
	return sql.Open(Driver, DSN)
}
