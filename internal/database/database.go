// Package database contains our database connection. This file hosts all the basics for
// creating a connection, wihle other files have the query functions.
package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/vinovest/sqlx"
)

type DB struct {
	conn *sqlx.DB
}

func New(databaseURI string) (*DB, error) {
	conn, err := sqlx.Connect("mysql", databaseURI)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn,
	}, nil
}
