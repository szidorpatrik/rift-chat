-- +goose Up
CREATE TABLE channels(
    id uuid NOT NULL DEFAULT uuidv7(),
    server_id uuid NOT NULL,
    name varchar(100) NOT NULL,
    "type" integer NOT NULL DEFAULT 0,
    position integer NOT NULL DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY(id),
    CONSTRAINT channels_server_id_fkey FOREIGN key(server_id) REFERENCES servers(id)
);

-- +goose Down
DROP TABLE channels;
