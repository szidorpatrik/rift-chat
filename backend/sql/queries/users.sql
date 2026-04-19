-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByUsernameOrEmail :one
SELECT * FROM users WHERE username = $1 OR email = $2;

-- name: CreateUser :one
INSERT INTO
    users (
        id,
        username,
        email,
        password_hash
    )
VALUES ($1, $2, $3, $4)
RETURNING
    *;
