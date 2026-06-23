// Package database contains our database connection. This file hosts all the basics for
// creating a connection, wihle other files have the query functions.
package database

import (
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vinovest/sqlx"
)

type DB struct {
	conn *sqlx.DB
}

func New(databaseURI string) (*DB, error) {
	dsn, err := URLtoDSN(databaseURI)
	if err != nil {
		return nil, err
	}

	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn,
	}, nil
}

func URLtoDSN(databaseURL string) (string, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}
	user := u.User.Username()
	pass, _ := u.User.Password()
	host := u.Host
	dbname := strings.TrimPrefix(u.Path, "/")

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname), nil
}
