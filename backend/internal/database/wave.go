package database

import (
	"context"
	"time"
)

type Wave struct {
	ID        int       `db:"id" json:"id"`
	State     string    `db:"state" json:"state"`
	BeginAt   time.Time `db:"begin_at" json:"begin_at"`
	CloseAt   time.Time `db:"close_at" json:"close_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (db *DB) GetWaveByID(ctx context.Context, id int) (*Wave, error) {
	var wave Wave
	if err := db.conn.GetContext(ctx, &wave, "SELECT id, state, begin_at, close_at, created_at FROM waves WHERE id = ?", id); err != nil {
		return nil, err
	}

	return &wave, nil
}
