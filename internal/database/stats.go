package database

import (
	"context"
)

type TraineeStat struct {
	ID              string `db:"id" json:"id"`
	Username        string `db:"username" json:"username"`
	Snowflake       string `db:"snowflake" json:"snowflake"`
	PublicMessages  int    `db:"public_messages" json:"public_messages"`
	PrivateMessages int    `db:"private_messages" json:"private_messages"`
	Cases           int    `db:"cases" json:"cases"`
	ThreadChat      int    `db:"thread_chat" json:"thread_chat"`
	ThreadReplies   int    `db:"thread_replies" json:"thread_replies"`
	ThreadClosures  int    `db:"thread_closures" json:"thread_closures"`
	SnippetsUsed    int    `db:"snippets_used" json:"snippets_used"`
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

	if err := db.conn.SelectContext(ctx, &trainees, `SELECT
    st.id,
    st.snowflake,
	st.username,
    COALESCE(SUM(sp.public_messages), 0)  AS public_messages,
    COALESCE(SUM(sp.private_messages), 0) AS private_messages,
    COALESCE(SUM(sp.cases), 0)            AS cases,
    COALESCE(SUM(sp.thread_chat), 0)      AS thread_chat,
    COALESCE(SUM(sp.thread_replies), 0)   AS thread_replies,
    COALESCE(SUM(sp.thread_closures), 0)  AS thread_closures,
    COALESCE(SUM(sp.snippets_used), 0)    AS snippets_used
FROM staff st
LEFT JOIN stats_per_date sp ON sp.user_id = st.id
WHERE st.wave_id = ?
GROUP BY st.id, st.snowflake
ORDER BY st.id`, waveID); err != nil {
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
