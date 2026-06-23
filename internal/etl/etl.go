// Package etl contains our Extract, Transform, Load functions for taking ModMail / Athena data and importing it to be used in academy.
package etl

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vinovest/sqlx"
)

// Extract data from one source, transform it in Go
// insert it into the new database, complete?

// Data to collect:
// - Modmail threads
// - Modail thread messages
// - Athena cases
// - Athena case notes
// during transform I want to calculate some stats. i want to work out when the pipeline was run and save
// that. we probably want to run this in a loop every 2 minutes, getting the current time at the start, and selecting all changes that have happened since then
// im guessing for threads we'd just do a "select all threads with trainees participating", CREATE or UPDATE + set new stats?
// for thread messages, cases, and case notes, i can just get created_at > start time? there needs to be an initial import that I can trigger though. preferably
// one that can be repeated

type Etl struct {
	startDate time.Time
	mmDB      *sqlx.DB
	outDB     *sqlx.DB
}

type ImportedThread struct {
	ID               BinaryUUID `db:"id"`
	Status           int        `db:"status"`
	UserID           string     `db:"user_id"`
	UserName         string     `db:"user_name"`
	Roles            string     `db:"roles"`
	CreatedAt        time.Time  `db:"created_at"`
	ClosedByID       *string    `db:"closed_by_id"`
	InboundMessages  int        `db:"inbound_messages"`
	OutboundMessages int        `db:"outbound_messages"`
	ChatMessages     int        `db:"chat_messages"`
}

type ImportedThreadMessage struct {
	ID          int        `db:"id"`
	ThreadID    BinaryUUID `db:"thread_id"`
	Kind        int        `db:"kind"`
	UserID      string     `db:"user_id"`
	UserName    string     `db:"user_name"`
	Body        string     `db:"body"`
	CreatedAt   time.Time  `db:"created_at"`
	Attachments string     `db:"attachments"`
	Metadata    string     `db:"metadata"`
}

func (e *Etl) OutTx() (*sqlx.Tx, error) {
	return e.outDB.Beginx()
}

func (e *Etl) FindAllTraineeThreads(ctx context.Context, traineeIDs []string) ([]ImportedThread, error) {
	query, args, err := sqlx.In(`
		SELECT
			t.id, t.status, t.user_id, t.user_name, coalesce(t.roles, '') roles, t.created_at, t.closed_by_id,
			COUNT(CASE WHEN tm.message_type = 3 THEN 1 END) AS inbound_messages,
			COUNT(CASE WHEN tm.message_type = 4 THEN 1 END) AS outbound_messages,
			COUNT(CASE WHEN tm.message_type = 2 THEN 1 END) AS chat_messages
		FROM threads t
		INNER JOIN thread_messages tm ON tm.thread_id = t.id
		WHERE tm.user_id IN (?)
		AND (t.created_at > ? OR t.updated_at > ?)
		GROUP BY t.id, t.status, t.user_id, t.user_name, t.roles, t.created_at, t.closed_by_id`,
		traineeIDs, e.startDate, e.startDate,
	)
	if err != nil {
		return nil, fmt.Errorf("building FindAllTraineeThreads query: %w", err)
	}
	query = e.mmDB.Rebind(query)

	threads := []ImportedThread{}
	if err := e.mmDB.SelectContext(ctx, &threads, query, args...); err != nil {
		return nil, fmt.Errorf("FindAllTraineeThreads: %w", err)
	}
	return threads, nil
}

func (e *Etl) FindThreadMessages(ctx context.Context, threadID string) ([]ImportedThreadMessage, error) {
	messages := []ImportedThreadMessage{}
	if err := e.mmDB.SelectContext(ctx, &messages, `
		SELECT id, thread_id, message_type AS kind, user_id, user_name,
		       body, created_at, attachments, metadata
		FROM thread_messages
		WHERE thread_id = ?
		ORDER BY created_at ASC`,
		threadID,
	); err != nil {
		return nil, fmt.Errorf("FindThreadMessages(%s): %w", threadID, err)
	}
	return messages, nil
}

func (e *Etl) InsertImportedThread(ctx context.Context, tx *sqlx.Tx, thread ImportedThread) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO threads (
			id, status, user_id, user_name, created_at, closed_by_id, roles, inbound_messages, outbound_messages, chat_messages
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		) AS new_data
		ON DUPLICATE KEY UPDATE
			status       = new_data.status,
			closed_by_id = new_data.closed_by_id`,
		thread.ID, thread.Status, thread.UserID, thread.UserName,
		thread.CreatedAt, thread.ClosedByID, thread.Roles, thread.InboundMessages, thread.OutboundMessages, thread.ChatMessages,
	)
	if err != nil {
		return fmt.Errorf("InsertImportedThread(%s): %w", thread.ID, err)
	}
	return nil
}

func (e *Etl) InsertThreadMessage(ctx context.Context, tx *sqlx.Tx, message ImportedThreadMessage) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO thread_messages (
			id, thread_id, kind, user_id, user_name, body, created_at, attachments, metadata
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		) AS new_data
		ON DUPLICATE KEY UPDATE
			body        = new_data.body,
			attachments = new_data.attachments,
			metadata    = new_data.metadata`,
		message.ID, message.ThreadID, message.Kind, message.UserID,
		message.UserName, message.Body, message.CreatedAt,
		message.Attachments, message.Metadata,
	)
	if err != nil {
		return fmt.Errorf("InsertThreadMessage(%d): %w", message.ID, err)
	}
	return nil
}

// id BINARY(16) PRIMARY KEY,
// -- 1 = open, 2 = closed, 3 = suspended.
// status INT NOT NULL DEFAULT 0,
// user_id VARCHAR(22) NOT NULL,
// user_name VARCHAR(128) NOT NULL,
// created_at TIMESTAMP NOT NULL,
// imported_at TIMESTAMP NOT NULL DEFAULT NOW(),
// closed_by_id VARCHAR(22) NOT NULL,
// roles TEXT NOT NULL,
// -- The following are stats that are calculated every time messages are imported.
// inbound_messages INT DEFAULT 0,
// outbound_messages INT DEFAULT 0,
// chat_messages INT DEFAULT 0

// UUID type, neeed for scanning in and pushing out

type BinaryUUID uuid.UUID

// Value converts to BINARY(16) when writing to DB
func (b BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(b).MarshalBinary()
}

// Scan converts from BINARY(16) when reading from DB
func (b *BinaryUUID) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		if len(v) == 16 {
			// Already binary
			parsed, err := uuid.FromBytes(v)
			if err != nil {
				return fmt.Errorf("BinaryUUID: %w", err)
			}
			*b = BinaryUUID(parsed)
		} else {
			// VARCHAR coming back as []byte
			parsed, err := uuid.ParseBytes(v)
			if err != nil {
				return fmt.Errorf("BinaryUUID: %w", err)
			}
			*b = BinaryUUID(parsed)
		}
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("BinaryUUID: %w", err)
		}
		*b = BinaryUUID(parsed)
	default:
		return fmt.Errorf("BinaryUUID: expected []byte or string, got %T", src)
	}
	return nil
}

func (b BinaryUUID) String() string {
	return uuid.UUID(b).String()
}
