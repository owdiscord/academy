package database

import (
	"context"
	"fmt"

	"github.com/owdiscord/academy/internal/formatting"
)

type Case struct {
	ID               int        `db:"id" json:"id"`
	CaseNumber       int        `db:"case_number" json:"case_number,omitempty"`
	ModID            string     `db:"mod_id" json:"mod_id"`
	ModName          string     `db:"mod_name" json:"mod_name"`
	ActionedUserID   string     `db:"actioned_user_id" json:"actioned_user_id"`
	ActionedUserName string     `db:"actioned_user_name" json:"actioned_user_name"`
	CreatedAt        int        `db:"created_at" json:"created_at"`
	Kind             int        `db:"type" json:"type"`
	Notes            []CaseNote `db:"-" json:"notes,omitempty"`
}

type CaseNote struct {
	ID        int    `db:"id" json:"id"`
	CaseID    int    `db:"case_id" json:"case_id,omitempty"`
	ModID     string `db:"mod_id" json:"mod_id"`
	Body      string `db:"body" json:"body"`
	CreatedAt int    `db:"created_at" json:"created_at"`
}

func (db *DB) GetAllCases(ctx context.Context, page, limit int) ([]Case, error) {
	cases := []Case{}

	if err := db.conn.SelectContext(ctx, &cases, "SELECT c.id, c.case_number, c.mod_id, COALESCE(s.username, 'Unknown') mod_name, c.actioned_user_id, c.actioned_user_name, UNIX_TIMESTAMP(c.created_at) created_at, c.type FROM cases c LEFT JOIN staff s ON s.snowflake = c.mod_id ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, (page-1)*limit); err != nil {
		return nil, err
	}

	return cases, nil
}

func (db *DB) GetCaseByID(ctx context.Context, id int) (*Case, error) {
	modCase := Case{}

	if err := db.conn.GetContext(ctx, &modCase, "SELECT c.id, c.case_number, c.mod_id, COALESCE(s.username, 'Unknown') mod_name, c.actioned_user_id, c.actioned_user_name, UNIX_TIMESTAMP(c.created_at) created_at, c.type FROM cases c LEFT JOIN staff s ON s.snowflake = c.mod_id WHERE c.id = ?", id); err != nil {
		return nil, fmt.Errorf("case(%d): %v", id, err)
	}

	notes := []CaseNote{}
	if err := db.conn.SelectContext(ctx, &notes, "SELECT id, body, mod_id, UNIX_TIMESTAMP(created_at) created_at FROM case_notes WHERE case_id = ? ORDER BY created_at ASC", id); err != nil {
		return nil, fmt.Errorf("case(%d) notes: %v", id, err)
	}

	for _, note := range notes {
		note.Body = string(formatting.MDtoHTML([]byte(note.Body)))

		modCase.Notes = append(modCase.Notes, note)
	}

	return &modCase, nil
}
