-- +goose Up
CREATE TABLE messages(
    id uuid NOT NULL DEFAULT uuidv7(),
    channel_id uuid NOT NULL,
    user_id uuid NOT NULL,
    content text NOT NULL,
    -- JSONB for flexible metadata (attachments, reactions, embeds, edited history)
    -- Example: {"reactions": {"thumbsup": 5}, "attachments": ["url1.jpg"]}
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    PRIMARY KEY(id),
    CONSTRAINT messages_channel_id_fkey FOREIGN key(channel_id) REFERENCES channels(id),
    CONSTRAINT messages_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE messages;