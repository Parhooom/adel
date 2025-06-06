-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS testcases (
    id BIGSERIAL PRIMARY KEY,
    problem_id BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input_data TEXT NOT NULL,
    output_data TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_testcases_problem_id ON testcases(problem_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_testcases_problem_id;
DROP TABLE testcases;
-- +goose StatementEnd
