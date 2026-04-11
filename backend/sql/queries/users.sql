-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;
