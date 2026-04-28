-- +goose Up
ALTER TABLE servers
ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE servers ADD COLUMN description TEXT;

ALTER TABLE servers
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

ALTER TABLE servers ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE servers DROP COLUMN description;

ALTER TABLE servers DROP COLUMN is_public;

ALTER TABLE servers DROP COLUMN updated_at;

ALTER TABLE servers DROP COLUMN deleted_at;