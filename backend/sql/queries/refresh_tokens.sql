-- name: CreateRefreshToken :one
INSERT INTO
    refresh_tokens (
        user_id,
        token,
        created_at,
        updated_at,
        expires_at,
        revoked_at
    )
VALUES (
        $1,
        $2,
        NOW(),
        NOW(),
        NOW() + INTERVAL '30 days',
        NULL
    )
RETURNING
    *;

-- name: GetRefreshTokensByUserId :many
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE
    token = $1;