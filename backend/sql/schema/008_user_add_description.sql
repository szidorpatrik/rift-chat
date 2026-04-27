-- +goose Up
ALTER TABLE users
ADD COLUMN description TEXT;

-- +goose Down
ALTER TABLE users
DROP COLUMN description;
