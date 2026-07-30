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

-- name: GetOwnerIdByServerId :one
SELECT owner_id FROM servers WHERE id = $1;

-- name: GetServerMembers :many
SELECT *
FROM users
    INNER JOIN server_members ON server_members.user_id = users.id
WHERE
    server_members.server_id = $1;

-- name: GetServerMemberIds :many
SELECT user_id FROM server_members WHERE server_id = $1;

-- name: LeaveServer :one
WITH
    server_info AS (
        SELECT owner_id
        FROM servers
        WHERE
            id = $1
    ),
    deleted AS (
        DELETE FROM server_members
        WHERE
            server_id = $1
            AND user_id = $2
            AND NOT EXISTS (
                SELECT 1
                FROM server_info
                WHERE
                    owner_id = $2
            )
        RETURNING
            user_id
    )
SELECT COALESCE(
        (
            SELECT owner_id = $2
            FROM server_info
        ), false
    )::boolean AS is_owner, EXISTS (
        SELECT 1
        FROM deleted
    )::boolean AS was_deleted;