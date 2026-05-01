-- +goose Up
CREATE TABLE server_members(
    server_id uuid NOT NULL,
    user_id uuid NOT NULL,
    server_name varchar(32),
    joined_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY(server_id,user_id),
    CONSTRAINT server_members_server_id_fkey FOREIGN key(server_id) REFERENCES servers(id),
    CONSTRAINT server_members_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE server_members;
