package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

const expiry = 7 * 24 * time.Hour

type Session struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	WaveID    int       `db:"wave_id" json:"wave_id"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Role      string    `db:"role" json:"role"`
}

func (db *DB) CreateSession(ctx context.Context, userID int, waveID int) (string, time.Time, error) {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	expires := time.Now().Add(expiry)

	_, err := db.conn.ExecContext(ctx, `INSERT INTO sessions (
		token, user_id, wave_id, expires_at,
	) VALUES (?, ?, ?, ?)`, token, userID, waveID, expires)

	return token, expires, err
}

func (db *DB) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	var session Session
	if err := db.conn.GetContext(ctx, &session, `SELECT 
		s.user_id, s.wave_id, s.expires_at, u.role 
    	FROM sessions s 
		INNER JOIN staff u ON u.id = s.user_id WHERE s.token = ? AND expires_at > NOW()`, token); err != nil {
		return nil, err
	}

	return &session, nil
}

func (db *DB) DeleteSessionByToken(ctx context.Context, token string) error {
	if _, err := db.conn.ExecContext(ctx, `DELETE FROM sessions WHERE token = ?`, token); err != nil {
		return err
	}

	return nil
}
