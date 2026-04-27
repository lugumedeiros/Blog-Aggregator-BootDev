-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: RemoveUser :exec
DELETE FROM users WHERE name = $1;

-- name: Reset :exec
DELETE FROM users;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;