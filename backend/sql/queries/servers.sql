-- name: CreateServer :one
INSERT INTO
    servers (id, name, owner_id)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: JoinServer :exec
INSERT INTO server_members (server_id, user_id) VALUES ($1, $2);