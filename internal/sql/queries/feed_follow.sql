-- name: CreateFeedFollow :exec
INSERT INTO feed_follow(created_at, updated_at, user_id, feed_id)
VALUES (
    $1, $2, $3, $4
);

-- name: GetFeedsByUser :many
SELECT * FROM feed_follow
WHERE user_id = $1;

-- name: GetFeedsByFeed :many
SELECT * FROM feed_follow
WHERE feed_id = $1;

-- name: GetUsersByFeed :many
SELECT user_id FROM feed_follow WHERE feed_id = $1;

-- name: GetFeedByUser :many
SELECT feed_id FROM feed_follow WHERE user_id = $1;