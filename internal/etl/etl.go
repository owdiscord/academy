// Package etl contains our Extract, Transform, Load functions for taking ModMail / Athena data and importing it to be used in academy.
package etl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/owdiscord/academy/internal/database"
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
	StartDate       time.Time
	WaveID          int
	StaffIDs        []string
	privateChannels []string
	statCollection  map[string]*DateStatParams
	athDB           *sqlx.DB
	mmDB            *sqlx.DB
	outDB           *sqlx.DB
}

func New(waveID int, from time.Time, athenaDB *sqlx.DB, modmailDB *sqlx.DB, outDB *sqlx.DB, staff []database.Staff, privateChannels []string) *Etl {
	staffIDs := []string{}
	statCollection := map[string]*DateStatParams{}

	for _, member := range staff {
		staffIDs = append(staffIDs, member.Snowflake)
		statCollection[member.Snowflake] = &DateStatParams{
			WaveID:         waveID,
			UserID:         member.ID,
			PublicMsgs:     0,
			PrivateMsgs:    0,
			Cases:          0,
			ThreadChat:     0,
			ThreadReplies:  0,
			ThreadClosures: 0,
			SnippetsUsed:   0,
		}
	}

	return &Etl{
		StartDate:       from,
		WaveID:          waveID,
		StaffIDs:        staffIDs,
		privateChannels: privateChannels,
		statCollection:  statCollection,
		athDB:           athenaDB,
		mmDB:            modmailDB,
		outDB:           outDB,
	}
}

func (e *Etl) IncreasePublicMsgStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].PublicMsgs += by
	}
}

func (e *Etl) IncreasePrivateMsgStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].PrivateMsgs += by
	}
}

func (e *Etl) IncreaseCasesStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].Cases += by
	}
}

func (e *Etl) IncreaseThreadChatStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].ThreadChat += by
	}
}

func (e *Etl) IncreaseThreadReplyStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].ThreadReplies += by
	}
}

func (e *Etl) IncreaseCloseStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].ThreadClosures += by
	}
}

func (e *Etl) IncreaseSnippetStat(snowflake string, by int) {
	if e.statCollection[snowflake] != nil {
		e.statCollection[snowflake].SnippetsUsed += by
	}
}

func (e *Etl) OutTx() (*sqlx.Tx, error) {
	return e.outDB.Beginx()
}

//
// # Threads
//

type ImportedThread struct {
	ID               database.BinaryUUID `db:"id"`
	Status           int                 `db:"status"`
	UserID           string              `db:"user_id"`
	UserName         string              `db:"user_name"`
	Roles            string              `db:"roles"`
	CreatedAt        time.Time           `db:"created_at"`
	ClosedByID       *string             `db:"closed_by_id"`
	InboundMessages  int                 `db:"inbound_messages"`
	OutboundMessages int                 `db:"outbound_messages"`
	ChatMessages     int                 `db:"chat_messages"`
}

type ImportedThreadMessage struct {
	ID          int                 `db:"id"`
	ThreadID    database.BinaryUUID `db:"thread_id"`
	Kind        int                 `db:"kind"`
	Role        string              `db:"role_name"`
	Anonymous   bool                `db:"is_anonymous"`
	UserID      string              `db:"user_id"`
	UserName    string              `db:"user_name"`
	Body        string              `db:"body"`
	CreatedAt   time.Time           `db:"created_at"`
	Attachments string              `db:"attachments"`
	Metadata    string              `db:"metadata"`
}

func (e *Etl) FindAllTraineeThreads(ctx context.Context) ([]ImportedThread, error) {
	query, args, err := sqlx.In(`
		SELECT
			t.id, t.status, t.user_id, t.user_name, coalesce(t.roles, '[]') roles, t.created_at, t.closed_by_id,
			COUNT(CASE WHEN tm.message_type = 3 THEN 1 END) AS inbound_messages,
			COUNT(CASE WHEN tm.message_type = 4 THEN 1 END) AS outbound_messages,
			COUNT(CASE WHEN tm.message_type = 2 THEN 1 END) AS chat_messages
		FROM threads t
		INNER JOIN thread_messages tm ON tm.thread_id = t.id
		WHERE tm.user_id IN (?)
		AND (t.created_at > ? OR t.updated_at > ?)
		GROUP BY t.id, t.status, t.user_id, t.user_name, t.roles, t.created_at, t.closed_by_id`,
		e.StaffIDs, e.StartDate, e.StartDate,
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
		       body, created_at, COALESCE(attachments, '[]') attachments, COALESCE(metadata, '{}') metadata, is_anonymous, COALESCE(role_name, 'Unknown') role_name
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
			id, status, user_id, user_name, created_at, closed_by_id, roles, inbound_messages, outbound_messages, chat_messages, wave_id, participants
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '[]'
		) AS new_data
		ON DUPLICATE KEY UPDATE
			status       = new_data.status,
			closed_by_id = new_data.closed_by_id`,
		thread.ID, thread.Status, thread.UserID, thread.UserName,
		thread.CreatedAt, thread.ClosedByID, thread.Roles, thread.InboundMessages, thread.OutboundMessages, thread.ChatMessages, e.WaveID,
	)
	if err != nil {
		return fmt.Errorf("InsertImportedThread(%s): %w", thread.ID, err)
	}
	return nil
}

func (e *Etl) InsertThreadMessage(ctx context.Context, tx *sqlx.Tx, message ImportedThreadMessage) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO thread_messages (
			id, thread_id, kind, user_id, user_name, body, created_at, attachments, metadata, anonymous, role
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		) AS new_data
		ON DUPLICATE KEY UPDATE
			body        = new_data.body,
			attachments = new_data.attachments,
			metadata    = new_data.metadata`,
		message.ID, message.ThreadID, message.Kind, message.UserID,
		message.UserName, message.Body, message.CreatedAt,
		message.Attachments, message.Metadata, message.Anonymous, message.Role,
	)
	if err != nil {
		return fmt.Errorf("InsertThreadMessage(%d): %w", message.ID, err)
	}
	return nil
}

func (e *Etl) RecalculateThreadMessageCounts(ctx context.Context, tx *sqlx.Tx, threadID database.BinaryUUID) error {
	_, err := tx.ExecContext(ctx, `
    UPDATE threads t SET
        participants      = COALESCE((SELECT JSON_ARRAYAGG(user_id) FROM (SELECT DISTINCT user_id FROM thread_messages WHERE thread_id = ? AND kind IN (2, 3)) AS u), '[]'),
        inbound_messages  = (SELECT COUNT(*) FROM thread_messages WHERE thread_id = t.id AND kind = 1),
        outbound_messages = (SELECT COUNT(*) FROM thread_messages WHERE thread_id = t.id AND kind = 2),
        chat_messages     = (SELECT COUNT(*) FROM thread_messages WHERE thread_id = t.id AND kind = 3)
    WHERE id = ?`,
		threadID, threadID,
	)
	if err != nil {
		return fmt.Errorf("RecalculateThreadMessageCounts(%s): %w", threadID, err)
	}
	return nil
}

//
// # Cases
//

type ImportedCase struct {
	ID         uint      `db:"id"`
	CaseNumber uint      `db:"case_number"`
	UserID     uint64    `db:"user_id"`
	UserName   string    `db:"user_name"`
	ModID      *uint64   `db:"mod_id"`
	Type       uint      `db:"type"`
	CreatedAt  time.Time `db:"created_at"`
	IsHidden   uint8     `db:"is_hidden"`
}

type ImportedCaseNote struct {
	ID        uint      `db:"id"`
	CaseID    uint      `db:"case_id"`
	ModID     *uint64   `db:"mod_id"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

func (e *Etl) FindAllTraineeCases(ctx context.Context) ([]ImportedCase, error) {
	query, args, err := sqlx.In(`
        SELECT
            id, case_number, user_id, user_name,
            mod_id, type, created_at
        FROM cases
        WHERE mod_id IN (?)
        AND created_at > ?`,
		e.StaffIDs, e.StartDate,
	)
	if err != nil {
		return nil, fmt.Errorf("building FindAllTraineeCases query: %w", err)
	}
	query = e.athDB.Rebind(query)

	cases := []ImportedCase{}
	if err := e.athDB.SelectContext(ctx, &cases, query, args...); err != nil {
		return nil, fmt.Errorf("FindAllTraineeCases: %w", err)
	}
	return cases, nil
}

func (e *Etl) FindCaseNotes(ctx context.Context, caseID uint) ([]ImportedCaseNote, error) {
	notes := []ImportedCaseNote{}
	if err := e.athDB.SelectContext(ctx, &notes, `
        SELECT id, case_id, mod_id, body, created_at
        FROM case_notes
        WHERE case_id = ?`,
		caseID,
	); err != nil {
		return nil, fmt.Errorf("FindCaseNotes(%d): %w", caseID, err)
	}
	return notes, nil
}

func (e *Etl) InsertImportedCase(ctx context.Context, tx *sqlx.Tx, c ImportedCase) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO cases (
            id, case_number, actioned_user_id, actioned_user_name,
            mod_id, type, created_at, wave_id
        ) VALUES (
            ?, ?, ?, ?,
			?, ?, ?, ?
        ) AS new_data
        ON DUPLICATE KEY UPDATE
            mod_id        = new_data.mod_id`,
		c.ID, c.CaseNumber, c.UserID, c.UserName,
		c.ModID, c.Type, c.CreatedAt, e.WaveID,
	)
	if err != nil {
		return fmt.Errorf("InsertImportedCase(%d): %w", c.ID, err)
	}
	return nil
}

func (e *Etl) InsertCaseNote(ctx context.Context, tx *sqlx.Tx, note ImportedCaseNote) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO case_notes (
            id, case_id, mod_id, body, created_at
        ) VALUES (
            ?, ?, ?, ?, ?
        ) AS new_data
        ON DUPLICATE KEY UPDATE
            body     = new_data.body`,
		note.ID, note.CaseID, note.ModID, note.Body, note.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("InsertCaseNote(%d): %w", note.ID, err)
	}
	return nil
}

//
// # Message stats
//

type MessageStat struct {
	UserID  string `db:"user_id"`
	Private int    `db:"private_messages"`
	Public  int    `db:"public_messages"`
}

func (e *Etl) GetMessageStats(ctx context.Context, tx *sqlx.Tx) ([]MessageStat, error) {
	chanIDs := strings.Join(e.privateChannels, ", ")
	userIDs := strings.Join(e.StaffIDs, ", ")

	stats := []MessageStat{}

	// You may think it's bad form to concatenate strings for a query.
	// You would be 100% correct, but I'm gonna save myself a heck of a lot of time by doing this lol
	if err := e.athDB.SelectContext(ctx, &stats, `
SELECT
    user_id,
    SUM(CASE WHEN channel_id IN (`+chanIDs+`) THEN 1 ELSE 0 END) AS private_messages,
    SUM(CASE WHEN channel_id NOT IN (`+chanIDs+`) THEN 1 ELSE 0 END) AS public_messages
FROM messages
WHERE posted_at > ? 
	AND user_id IN (`+userIDs+`)
  AND is_bot = 0
GROUP BY user_id`, e.StartDate); err != nil {
		return nil, err
	}

	return stats, nil
}

type DateStatParams struct {
	WaveID         int
	UserID         int
	PublicMsgs     int
	PrivateMsgs    int
	Cases          int
	ThreadChat     int
	ThreadReplies  int
	ThreadClosures int
	SnippetsUsed   int
}

// SaveDateStatsForUser is a mega-function for updating every value, which should be seldom-used
func (e *Etl) SaveDateStatsForUser(ctx context.Context, tx *sqlx.Tx, params DateStatParams) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO stats_per_date
  (date, user_id, wave_id, public_messages, private_messages, cases, thread_chat, thread_replies, thread_closures, snippets_used)
VALUES
  (CURRENT_DATE, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  wave_id            = VALUES(wave_id),
  public_messages    = public_messages    + VALUES(public_messages),
  private_messages   = private_messages   + VALUES(private_messages),
  cases              = cases              + VALUES(cases),
  thread_chat        = thread_chat        + VALUES(thread_chat),
  thread_replies     = thread_replies     + VALUES(thread_replies),
  thread_closures    = thread_closures    + VALUES(thread_closures),
  snippets_used      = snippets_used      + VALUES(snippets_used)`, params.UserID, params.WaveID, params.PublicMsgs, params.PrivateMsgs, params.Cases, params.ThreadChat, params.ThreadReplies, params.ThreadClosures, params.SnippetsUsed)

	return err
}

func (e *Etl) SaveAllDateStats(ctx context.Context, tx *sqlx.Tx) error {
	for _, stats := range e.statCollection {
		if err := e.SaveDateStatsForUser(ctx, tx, *stats); err != nil {
			return err
		}
	}

	return nil
}
