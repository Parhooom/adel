-- +goose Up
-- +goose StatementBegin
ALTER TABLE testcases ADD COLUMN user_id BIGINT REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE testcases DROP COLUMN user_id;
-- +goose StatementEnd