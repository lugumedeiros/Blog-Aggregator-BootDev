-- name: CreateFeed :exec
INSERT INTO feeds(user_id, name, url)
VALUES(
    $1,
    $2,
    $3
);

-- name: GetFeeds :many
SELECT * FROM feeds;