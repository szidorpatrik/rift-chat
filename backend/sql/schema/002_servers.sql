-- +goose Up
CREATE TABLE servers(
    id uuid NOT NULL DEFAULT uuidv7(),
    name varchar(100) NOT NULL,
    owner_id uuid NOT NULL,
    icon_url varchar(255),
    created_at timestamp with time zone DEFAULT now(),
    is_public boolean NOT NULL DEFAULT false,
    description text,
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    background_url varchar(255),
    PRIMARY KEY(id),
    CONSTRAINT servers_owner_id_fkey FOREIGN key(owner_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE servers;
