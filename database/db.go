package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// package database contains the stuff we need to interact with the underlying postgres database

// OpenDB opens a database and returns it for other functions to use
func OpenDB() (*sql.DB, error) {
	connStr := "user=ghost dbname=openclimate sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, errors.Wrap(err, "could not open database, quitting")
	}
	return db, nil
}

