-- +goose Up
ALTER TABLE servers
ADD COLUMN background_url TEXT;

-- +goose Down
ALTER TABLE servers
DROP COLUMN background_url;
