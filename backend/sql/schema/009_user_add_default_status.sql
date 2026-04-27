-- +goose Up
ALTER TABLE users
ALTER COLUMN status SET DEFAULT 'offline';

-- +goose Down
ALTER TABLE users
ALTER COLUMN status DROP DEFAULT;
