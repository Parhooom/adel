-- +goose Up
-- +goose StatementBegin
ALTER TABLE submissions ADD COLUMN user_id BIGINT REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE submissions DROP COLUMN user_id;
-- +goose StatementEnd