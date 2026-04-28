-- +goose Up
ALTER TABLE servers
ALTER COLUMN name TYPE TEXT;

-- +goose Down
ALTER TABLE servers
ALTER COLUMN name TYPE VARCHAR(100);
