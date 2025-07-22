package postgres

func (pg *PostgresSubmissionStore) GetSubmissionsByUserID(userID int64) ([]Submission, error) {
	query := `
	SELECT id, user_id, problem_id, code, language, status, 
			execution_time_ms, memory_usage_mb, error_message, 
			created_at, updated_at
	FROM submissions 
	WHERE user_id = $1 
	ORDER BY created_at DESC
    `

	rows, err := pg.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var submission Submission
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ProblemID,
			&submission.Code,
			&submission.Language,
			&submission.Status,
			&submission.ExecutionTimeMs,
			&submission.MemoryUsageMB,
			&submission.ErrorMessage,
			&submission.CreatedAt,
			&submission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return submissions, nil
}

func (pg *PostgresSubmissionStore) GetUserSolvedProblemsCount(userID int64) (int, error) {
	query := `
	SELECT COUNT(DISTINCT problem_id)
	FROM submissions 
	WHERE user_id = $1 AND status = 'accepted'
	`

	var count int
	err := pg.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pg *PostgresSubmissionStore) GetUserSuccessRate(userID int64) (float64, error) {
	query := `
	SELECT 
		COUNT(*) FILTER (WHERE status = 'accepted') as accepted_count,
		COUNT(*) as total_count
	FROM submissions 
	WHERE user_id = $1
	`

	var acceptedCount, totalCount int
	err := pg.db.QueryRow(query, userID).Scan(&acceptedCount, &totalCount)
	if err != nil {
		return 0, err
	}

	if totalCount == 0 {
		return 0, nil
	}

	successRate := float64(acceptedCount) / float64(totalCount) * 100
	return successRate, nil
}

func (pg *PostgresSubmissionStore) GetTotalSubmissionsCount() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM submissions`

	err := pg.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *PostgresUserStore) GetTotalUsersCount() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM users`

	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
