package database

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

type SQLDatabase struct {
	database *sql.DB
}

func NewDatabase(db *sql.DB) SQLDatabase {
	return SQLDatabase{
		database: db,
	}
}

var ErrNotFound = errors.New("not found in db")

func NewPostgresConn() SQLDatabase {
	uri, found := os.LookupEnv("POSTGRES_URI")
	if !found {
		panic("POSTGRES_URI not found")
	}
	db, err := sql.Open("pgx", uri)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return NewDatabase(db)
}
