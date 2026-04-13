-- name: CreateMessage :one
INSERT INTO
    messages (
        id,
        channel_id,
        user_id,
        content,
        metadata
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetMessagesForChannel :many
SELECT *
FROM messages
WHERE
    channel_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;