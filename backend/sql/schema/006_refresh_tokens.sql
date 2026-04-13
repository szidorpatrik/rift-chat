-- +goose Up
CREATE TABLE refresh_tokens (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE refresh_tokens;
