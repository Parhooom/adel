package postgres

import (
	"database/sql"
	"time"
)

type Submission struct {
	ID              int64     `json:"id"`
	ProblemID       int64     `json:"problem_id"`
	Code            string    `json:"code"`
	Language        string    `json:"language"`
	Status          string    `json:"status"`
	ExecutionTimeMs int64     `json:"execution_time_ms"`
	MemoryUsageMB   int64     `json:"memory_usage_mb"`
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PostgresSubmissionStore struct {
	db *sql.DB
}

func NewPostgresSubmissionStore(db *sql.DB) *PostgresSubmissionStore {
	return &PostgresSubmissionStore{db: db}
}

type SubmissionStore interface {
	CreateSubmission(submission *Submission) (*Submission, error)
	GetSubmissionByID(id int64) (*Submission, error)
	DeleteSubmission(id int64) error
	UpdateSubmission(submission *Submission) error
}

func (pg *PostgresSubmissionStore) CreateSubmission(submission *Submission) (*Submission, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO submissions (problem_id, code, language, status, execution_time_ms, memory_usage_mb, error_message)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	err = tx.QueryRow(query, submission.ProblemID, submission.Code, submission.Language, submission.Status, submission.ExecutionTimeMs, submission.MemoryUsageMB, submission.ErrorMessage).Scan(&submission.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (pg *PostgresSubmissionStore) GetSubmissionByID(id int64) (*Submission, error) {
	submission := &Submission{}

	query := `
	SELECT id, problem_id, code, language, status, execution_time_ms, memory_usage_mb, error_message, created_at, updated_at
	FROM submissions
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&submission.ID, &submission.ProblemID, &submission.Code, &submission.Language, &submission.Status, &submission.ExecutionTimeMs, &submission.MemoryUsageMB, &submission.ErrorMessage, &submission.CreatedAt, &submission.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (pg *PostgresSubmissionStore) DeleteSubmission(id int64) error {
	query := `
	DELETE FROM submissions
	WHERE id = $1
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresSubmissionStore) UpdateSubmission(submission *Submission) error {
	query := `
	UPDATE submissions
	SET status = $1, execution_time_ms = $2, memory_usage_mb = $3, error_message = $4
	WHERE id = $5
	`
	result, err := pg.db.Exec(query, submission.Status, submission.ExecutionTimeMs, submission.MemoryUsageMB, submission.ErrorMessage, submission.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
