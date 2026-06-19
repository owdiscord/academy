package database

import "context"

func (db *DB) GetQuestions(ctx context.Context) ([]string, error) {
	questions := []string{}
	if err := db.conn.SelectContext(ctx, &questions, "SELECT text FROM interview_questions"); err != nil {
		return nil, err
	}

	return questions, nil
}
