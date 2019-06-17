package database

import (
	"github.com/pkg/errors"
	"database/sql"
	_ "github.com/lib/pq"
)

func OpenDB() (*sql.DB, error) {
	connStr := "user=ghost dbname=openclimate sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, errors.Wrap(err, "could not open database, quitting")
	}
	return db, nil
}
