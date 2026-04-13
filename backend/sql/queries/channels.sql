-- name: CreateChannel :one
INSERT INTO
    channels (
        id,
        server_id,
        name,
        type,
        position
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetChannelsForServer :many
SELECT * FROM channels WHERE server_id = $1 ORDER BY position ASC;

-- name: GetChannelById :one
SELECT * FROM channels WHERE id = $1;
