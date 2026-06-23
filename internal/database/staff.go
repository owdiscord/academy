package database

import (
	"context"

	"github.com/owdiscord/academy/internal/discord"
)

type Staff struct {
	ID          int    `db:"id" json:"id,omitempty"`
	WaveID      int    `db:"wave_id" json:"wave_id,omitempty"`
	Snowflake   string `db:"snowflake" json:"snowflake"`
	Username    string `db:"username" json:"username"`
	DisplayName string `db:"display_name" json:"display_name"`
	Role        string `db:"role" json:"role,omitempty"`
}

func (db *DB) UpdateUserDetails(ctx context.Context, details discord.DiscordUser) error {
	if _, err := db.conn.ExecContext(ctx, "UPDATE staff SET username = ?, display_name = ? WHERE snowflake = ?", details.Username, details.GlobalName, details.ID); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetWaveTrainees(ctx context.Context, waveID int) ([]Staff, error) {
	staff := []Staff{}
	if err := db.conn.SelectContext(ctx, &staff, "SELECT snowflake, username, display_name FROM staff WHERE wave_id = ? AND role = 'trainee'", waveID); err != nil {
		return nil, err
	}

	return staff, nil
}

func (db *DB) GetStaffDetails(ctx context.Context, userID int, waveID int) (*Staff, error) {
	var staff Staff
	if err := db.conn.GetContext(ctx, &staff, "SELECT id, wave_id, snowflake, username, display_name, role FROM staff WHERE id = ? AND wave_id = ?", userID, waveID); err != nil {
		return nil, err
	}

	return &staff, nil
}

func (db *DB) LatestUserForDiscordID(ctx context.Context, snowflake string) (*Staff, error) {
	var staff Staff
	if err := db.conn.GetContext(ctx, &staff, "SELECT id, wave_id, snowflake, username, display_name, role FROM staff WHERE snowflake = ? ORDER BY wave_id LIMIT 1", snowflake); err != nil {
		return nil, err
	}

	return &staff, nil
}
