-- +goose Up
CREATE TABLE credentials (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE credentials;