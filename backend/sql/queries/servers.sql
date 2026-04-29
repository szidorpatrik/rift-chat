-- name: CreateServer :one
INSERT INTO
    servers (
        id,
        owner_id,
        name,
        description,
        is_public
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: JoinServer :exec
INSERT INTO server_members (server_id, user_id) VALUES ($1, $2);

-- name: GetServersByUserId :many
SELECT *
FROM servers
WHERE
    id IN (
        SELECT server_id
        FROM server_members
        WHERE
            user_id = $1
    );

-- name: GetPublicServers :many
SELECT * FROM servers WHERE is_public = true;

-- name: GetServerById :one
SELECT * FROM servers WHERE id = $1;
