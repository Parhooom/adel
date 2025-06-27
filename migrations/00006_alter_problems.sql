-- +goose Up
-- +goose StatementBegin
ALTER TABLE problems ADD COLUMN user_id BIGINT REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE problems DROP COLUMN user_id;
-- +goose StatementEnd