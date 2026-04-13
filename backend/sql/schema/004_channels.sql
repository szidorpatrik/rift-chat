-- +goose Up
CREATE TABLE channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid_v7(),
    server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'text', -- 'text' or 'voice' or 'category'
    position INTEGER NOT NULL DEFAULT 0, -- For ordering channels in the UI
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE channels;
