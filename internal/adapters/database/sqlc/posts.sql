-- name: CreatePost :one
INSERT INTO posts (
    feed_id,
    title,
    url,
    description,
    published_at,
    guid
)
VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (feed_id, url) DO NOTHING
RETURNING *;


