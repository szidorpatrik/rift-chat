-- +goose Up
ALTER TABLE channels ALTER COLUMN type DROP DEFAULT;

ALTER TABLE channels
ALTER COLUMN type TYPE INTEGER USING CASE
    WHEN type = 'text' THEN 0
    WHEN type = 'voice' THEN 1
    WHEN type = 'category' THEN 2
    ELSE 0
END,
ALTER COLUMN type
SET DEFAULT 0;

-- +goose Down
ALTER TABLE channels ALTER COLUMN type DROP DEFAULT;

ALTER TABLE channels
ALTER COLUMN type TYPE VARCHAR(20) USING CASE
    WHEN type = 0 THEN 'text'
    WHEN type = 1 THEN 'voice'
    WHEN type = 2 THEN 'category'
    ELSE 'text'
END,
ALTER COLUMN type
SET DEFAULT 'text';
