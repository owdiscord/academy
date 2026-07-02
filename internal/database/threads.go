package database

import (
	"context"

	"github.com/owdiscord/academy/internal/formatting"
)

type Thread struct {
	ID               BinaryUUID      `db:"id" json:"id"`
	Status           int             `db:"status" json:"status"`
	WaveID           int             `db:"wave_id" json:"wave_id,omitempty"`
	UserID           string          `db:"user_id" json:"user_id"`
	UserName         string          `db:"user_name" json:"user_name"`
	CreatedAt        int             `db:"created_at" json:"created_at"`
	ClosedByID       *string         `db:"closed_by_id" json:"closed_by_id"`
	Roles            JSONStringArray `db:"roles" json:"roles"`
	InboundMessages  int             `db:"inbound_messages" json:"inbound_messages"`
	OutboundMessages int             `db:"outbound_messages" json:"outbound_messages"`
	ChatMessages     int             `db:"chat_messages" json:"chat_messages"`
	Participants     JSONStringArray `db:"participants" json:"participants"`
	Messages         []ThreadMessage `db:"-" json:"messages,omitempty"`
}

type ThreadMessage struct {
	ID          int             `db:"id" json:"id"`
	ThreadID    BinaryUUID      `db:"thread_id" json:"thread_id"`
	Kind        int             `db:"kind" json:"kind"`
	UserID      string          `db:"user_id" json:"user_id"`
	UserName    string          `db:"user_name" json:"user_name"`
	Role        string          `db:"role" json:"role"`
	Anonymous   bool            `db:"anonymous" json:"anonymous"`
	Body        string          `db:"body" json:"body"`
	CreatedAt   int             `db:"created_at" json:"created_at"`
	Attachments JSONStringArray `db:"attachments" json:"attachments"`
	Metadata    JSONMap         `db:"metadata" json:"metadata"`
}

func (db *DB) GetAllThreads(ctx context.Context, page, limit int) ([]Thread, error) {
	threads := []Thread{}

	if err := db.conn.SelectContext(ctx, &threads, "SELECT id, status, wave_id, user_id, user_name, UNIX_TIMESTAMP(created_at) created_at, closed_by_id, roles, participants, inbound_messages, outbound_messages, chat_messages FROM threads LIMIT ? OFFSET ? ORDER BY created_at DESC", limit, (page-1)*limit); err != nil {
		return nil, err
	}

	return threads, nil
}

func (db *DB) GetThreadByID(ctx context.Context, id BinaryUUID) (*Thread, error) {
	thread := Thread{}

	if err := db.conn.GetContext(ctx, &thread, "SELECT id, status, wave_id, user_id, user_name, UNIX_TIMESTAMP(created_at) created_at, closed_by_id, roles, participants, inbound_messages, outbound_messages, chat_messages FROM threads WHERE id = ?", id); err != nil {
		return nil, err
	}

	messages := []ThreadMessage{}
	if err := db.conn.SelectContext(ctx, &messages, "SELECT id, thread_id, kind, user_id, user_name, COALESCE(role, 'system') role, body, UNIX_TIMESTAMP(created_at) created_at, COALESCE(attachments, '[]') attachments, COALESCE(metadata, '{}') metadata FROM thread_messages WHERE thread_id = ? ORDER BY created_at ASC", id); err != nil {
		return nil, err
	}

	for _, msg := range messages {
		msg.Body = string(formatting.MDtoHTML([]byte(msg.Body)))
		thread.Messages = append(thread.Messages, msg)
	}

	return &thread, nil
}
