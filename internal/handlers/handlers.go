// Package handlers contains our HTTP handlers and central context struct
package handlers

import "github.com/owdiscord/academy/internal/database"

type Handlers struct {
	db *database.DB
}

func New(db *database.DB) Handlers {
	return Handlers{
		db,
	}
}
