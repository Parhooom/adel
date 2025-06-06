-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS problems (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(50) CHECK (difficulty IN ('easy', 'medium', 'hard')),
    time_limit_ms INT NOT NULL DEFAULT 2000, -- 2 seconds
    memory_limit_mb INT NOT NULL DEFAULT 2,  -- 2 MB
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_problems_difficulty ON problems(difficulty);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_problems_difficulty;
DROP TABLE problems;
-- +goose StatementEnd