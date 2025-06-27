package postgres

import (
	"database/sql"
)

type TestCase struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	ProblemID  int64  `json:"problem_id"`
	IsActive   bool   `json:"is_active"`
	InputData  string `json:"input_data"`
	OutputData string `json:"output_data"`
}

type Problem struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	TimeLimit   int        `json:"time_limit_ms"`
	MemoryLimit int        `json:"memory_limit_mb"`
	IsActive    bool       `json:"is_active"`
	TestCases   []TestCase `json:"testcases"`
}

type ProblemSummary struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

type PostgresProblemStore struct {
	db *sql.DB
}

func NewPostgresProblemStore(db *sql.DB) *PostgresProblemStore {
	return &PostgresProblemStore{db: db}
}

type ProblemStore interface {
	CreateProblem(*Problem) (*Problem, error)
	GetProblemByID(id int64) (*Problem, error)
	DeleteProblem(id int64) error
	UpdateProblem(problem *Problem) error
	GetProblemOwner(problemID int64) (int64, error)
	GetAllProblems() ([]ProblemSummary, error)
}

func (pg *PostgresProblemStore) CreateProblem(problem *Problem) (*Problem, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO problems (user_id, title, description, difficulty, time_limit_ms, memory_limit_mb, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	err = tx.QueryRow(query, problem.UserID, problem.Title, problem.Description, problem.Difficulty, problem.TimeLimit, problem.MemoryLimit, problem.IsActive).Scan(&problem.ID)
	if err != nil {
		return nil, err
	}

	for i := range problem.TestCases {
		problem.TestCases[i].ProblemID = problem.ID

		query := `
		INSERT INTO testcases (user_id, problem_id, input_data, output_data, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`
		err = tx.QueryRow(query, problem.UserID, problem.ID, problem.TestCases[i].InputData, problem.TestCases[i].OutputData, problem.TestCases[i].IsActive).Scan(&problem.TestCases[i].ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return problem, nil
}

func (pg *PostgresProblemStore) GetProblemByID(id int64) (*Problem, error) {
	problem := &Problem{}

	query := `
	SELECT p.id, p.user_id, p.title, p.description, p.difficulty, p.time_limit_ms, p.memory_limit_mb, p.is_active
	FROM problems p
	WHERE p.id = $1 AND p.is_active = true
	`

	err := pg.db.QueryRow(query, id).Scan(&problem.ID, &problem.UserID, &problem.Title, &problem.Description, &problem.Difficulty, &problem.TimeLimit, &problem.MemoryLimit, &problem.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	testcaseQuery := `
	SELECT tc.id, tc.user_id, tc.problem_id, tc.input_data, tc.output_data, tc.is_active
	FROM testcases tc
	WHERE tc.problem_id = $1 AND tc.is_active = true
	`
	rows, err := pg.db.Query(testcaseQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var testCase TestCase
		err := rows.Scan(
			&testCase.ID,
			&testCase.UserID,
			&testCase.ProblemID,
			&testCase.InputData,
			&testCase.OutputData,
			&testCase.IsActive,
		)
		if err != nil {
			return nil, err
		}

		problem.TestCases = append(problem.TestCases, testCase)
	}

	return problem, nil
}

func (pg *PostgresProblemStore) DeleteProblem(id int64) error {
	query := `
  DELETE from problems
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

func (pg *PostgresProblemStore) UpdateProblem(problem *Problem) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	UPDATE problems
	SET title = $1, description = $2, difficulty = $3, time_limit_ms = $4, memory_limit_mb = $5, is_active = $6
	WHERE id = $7
	`

	result, err := tx.Exec(query, problem.Title, problem.Description, problem.Difficulty, problem.TimeLimit, problem.MemoryLimit, problem.IsActive, problem.ID)
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

	_, err = tx.Exec(`DELETE FROM testcases WHERE problem_id = $1`, problem.ID)
	if err != nil {
		return err
	}

	for i := range problem.TestCases {
		problem.TestCases[i].ProblemID = problem.ID

		query := `
		INSERT INTO testcases (user_id, problem_id, input_data, output_data, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`

		err = tx.QueryRow(query, problem.UserID, problem.ID, problem.TestCases[i].InputData, problem.TestCases[i].OutputData, problem.TestCases[i].IsActive).Scan(&problem.TestCases[i].ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *PostgresProblemStore) GetProblemOwner(problemID int64) (int64, error) {
	var userID int64

	query := `
  SELECT user_id
  FROM problems
  WHERE id = $1
  `

	err := pg.db.QueryRow(query, problemID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (pg *PostgresProblemStore) GetAllProblems() ([]ProblemSummary, error) {
	problems := []ProblemSummary{}

	query := `
	SELECT p.id, p.title, p.description, p.difficulty
	FROM problems p
	WHERE p.is_active = true
	`

	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var problem ProblemSummary
		err := rows.Scan(
			&problem.ID,
			&problem.Title,
			&problem.Description,
			&problem.Difficulty,
		)
		if err != nil {
			return nil, err
		}

		problems = append(problems, problem)
	}

	return problems, nil
}
