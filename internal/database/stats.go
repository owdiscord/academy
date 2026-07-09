package database

import (
	"context"
)

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

type StatsPerDate struct {
	Date            string `db:"date" json:"date"`
	PublicMessages  int    `db:"public_messages" json:"public_messages"`
	PrivateMessages int    `db:"private_messages" json:"private_messages"`
	ThreadChat      int    `db:"thread_chat" json:"thread_chat"`
	ThreadReplies   int    `db:"thread_replies" json:"thread_replies"`
	ThreadClosures  int    `db:"thread_closures" json:"thread_closures"`
	Cases           int    `db:"cases" json:"cases"`
	SnippetsUsed    int    `db:"snippets_used" json:"snippets_used"`
}

func (db *DB) GetStatsOverview(ctx context.Context, waveID int) (*StatsOverview, error) {
	trainees := []TraineeStat{}

	if err := db.conn.SelectContext(ctx, &trainees, "SELECT snowflake user_id, username, thread_participation_count, message_count, thread_count, case_count FROM staff WHERE role = 'trainee' AND wave_id = ?", waveID); err != nil {
		return nil, err
	}

	var overview StatsOverview

	if err := db.conn.GetContext(ctx, &overview, `WITH w AS (SELECT ? AS wave_id)
	SELECT
		(SELECT COUNT(*) FROM issues WHERE wave_id = w.wave_id) AS issue_count,
		(SELECT COUNT(*) FROM threads WHERE wave_id = w.wave_id) AS thread_count,
		(SELECT SUM(public_messages) + SUM(private_messages) FROM stats_per_date WHERE wave_id = w.wave_id) AS message_count,
		(SELECT COUNT(*) FROM cases WHERE wave_id = w.wave_id) AS case_count
	FROM w`, waveID); err != nil {
		return nil, err
	}

	overview.TraineeStats = trainees
	return &overview, nil
}

func (db *DB) GetStatsForStaff(ctx context.Context, staffID int) ([]StatsPerDate, error) {
	stats := []StatsPerDate{}

	if err := db.conn.SelectContext(ctx, &stats, "SELECT date, public_messages, private_messages, thread_chat, thread_replies, thread_closures, cases, snippets_used FROM stats_per_date WHERE user_id = ? ORDER BY date", staffID); err != nil {
		return nil, err
	}

	return stats, nil
}
