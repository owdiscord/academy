package database

import "context"

type Issue struct {
	ID        int    `db:"id" json:"id"`
	CreatedBy *int   `db:"created_by" json:"created_by"`
	CreatedAt int    `db:"created_at" json:"created_at"`
	TraineeID *int   `db:"trainee_id" json:"trainee_id"`
	ThreadID  *int   `db:"thread_id" json:"thread_id"`
	MessageID *int   `db:"message_id" json:"message_id"`
	Status    string `db:"status" json:"status"`
	Reason    string `db:"reason" json:"reason"`
	Category  string `db:"category" json:"category"`
}

func (db *DB) GetFullIssues(ctx context.Context, waveID int) ([]Issue, error) {
	issues := []Issue{}
	if err := db.conn.SelectContext(ctx, &issues, "SELECT id, UNIX_TIMESTAMP(created_at) created_at, created_by, trainee_id, thread_id, message_id, status, reason, category FROM issues WHERE wave_id = ? AND status != 'archived' ORDER BY created_at", waveID); err != nil {
		return nil, err
	}

	return issues, nil
}

type CreateIssueParams struct {
	WaveID    int    `json:"wave_id"`
	CreatedBy *int   `json:"created_by"`
	TraineeID *int   `json:"trainee_id"`
	ThreadID  *int   `json:"thread_id"`
	MessageID *int   `json:"message_id"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Category  string `json:"category"`
}

func (db *DB) CreateIssue(ctx context.Context, params CreateIssueParams) error {
	_, err := db.conn.ExecContext(ctx, `INSERT INTO issues (
		wave_id, created_by, trainee_id, thread_id, message_id, status, reason, category	
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?
	)`, params.WaveID, params.CreatedBy, params.TraineeID, params.ThreadID, params.MessageID, params.Status, params.Reason, params.Category)

	return err
}
