-- +goose Up
CREATE INDEX IF NOT EXISTS idx_users_created_at_id on users (created_at, id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_created_at_id