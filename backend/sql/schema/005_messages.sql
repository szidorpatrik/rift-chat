-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    
    -- JSONB for flexible metadata (attachments, reactions, embeds, edited history)
    -- Example: {"reactions": {"thumbsup": 5}, "attachments": ["url1.jpg"]}
    metadata JSONB DEFAULT '{}'::jsonb, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE messages;
