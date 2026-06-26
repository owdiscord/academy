package database

import "context"

type TraineeStat struct {
	Username           string `db:"username" json:"username"`
	UserID             string `db:"user_id" json:"user_id"`
	MessageCount       int    `db:"message_count" json:"message_count"`
	ParticipationCount int    `db:"thread_participation_count" json:"thread_participation_count"`
	ThreadCount        int    `db:"thread_count" json:"thread_count"`
	CaseCount          int    `db:"case_count" json:"case_count"`
}

type StatsOverview struct {
	CaseCount    int           `db:"case_count" json:"case_count"`
	ThreadCount  int           `db:"thread_count" json:"thread_count"`
	MessageCount int           `db:"message_count" json:"message_count"`
	IssueCount   int           `db:"issue_count" json:"issue_count"`
	TraineeStats []TraineeStat `db:"-" json:"trainee_stats"`
}

func (db *DB) GetStatsOverview(ctx context.Context, waveID int) (*StatsOverview, error) {
	trainees := []TraineeStat{}

	if err := db.conn.SelectContext(ctx, &trainees, "SELECT snowflake user_id, username, thread_participation_count, message_count, thread_count, case_count FROM staff WHERE role = 'trainee' AND wave_id = ?", waveID); err != nil {
		return nil, err
	}

	var overview StatsOverview

	if err := db.conn.GetContext(ctx, &overview, `SELECT
    (SELECT COUNT(*) FROM issues WHERE wave_id = ?) AS issue_count,
    COALESCE(SUM(message_count), 0) AS message_count,
    COALESCE(SUM(thread_count), 0)  AS thread_count,
    COALESCE(SUM(case_count), 0)    AS case_count
FROM staff WHERE role = 'trainee' AND wave_id = ?;`, waveID, waveID); err != nil {
		return nil, err
	}

	overview.TraineeStats = trainees
	return &overview, nil
}
