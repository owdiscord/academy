package database

import "context"

type Staff struct {
	ID          int    `db:"id" json:"id"`
	WaveID      int    `db:"wave_id" json:"wave_id"`
	Snowflake   string `db:"snowflake" json:"snowflake"`
	Username    string `db:"username" json:"username"`
	DisplayName string `db:"display_name" json:"display_name"`
	Role        string `db:"role" json:"role"`
}

func (db *DB) GetStaffDetails(ctx context.Context, userID int, waveID int) (*Staff, error) {
	var staff Staff
	if err := db.conn.GetContext(ctx, &staff, "SELECT id, wave_id, snowflake, username, display_name, role FROM staff WHERE id = ? AND wave_id = ?", userID, waveID); err != nil {
		return nil, err
	}

	return &staff, nil
}

// // Ensure a given discord ID (snowflake) has permission to access at least one
// // wave, returning the latest wave ID.
// export async function latestUserForDiscordID(
//   sql: DbQuery,
//   discord_id: string,
// ): Promise<{ id: number; wave_id: number; role: string } | null> {
//   const res =
//     await sql`SELECT id, wave_id, role FROM academy_staff WHERE snowflake = ${discord_id} ORDER BY wave_id DESC LIMIT 1`;
//
//   return (res[0] as { id: number; wave_id: number; role: string }) || null;
// }
