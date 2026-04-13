-- +goose Up
CREATE TABLE server_members (
    server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    server_name VARCHAR(32),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (server_id, user_id)
);
-- +goose Down
DROP TABLE server_members;
