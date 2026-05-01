-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE users(
    id uuid NOT NULL DEFAULT uuidv7(),
    username varchar(32) NOT NULL,
    email varchar(255) NOT NULL,
    avatar_url varchar(255),
    status varchar(20) DEFAULT 'offline'::character varying,
    password_hash varchar(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    description varchar(255),
    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX users_username_key ON public.users USING btree (username);
CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email);

-- +goose Down
DROP TABLE users;