-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    (SELECT COALESCE(MAX(id), 0) + 1 FROM users),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: RemoveUser :exec
DELETE FROM users WHERE name = $1;