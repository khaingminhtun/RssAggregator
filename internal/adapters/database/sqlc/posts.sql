-- name: CreatePost :one
INSERT INTO posts (feed_id, title, url, description, published_at, guid)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (feed_id, url) DO UPDATE
    SET title = EXCLUDED.title
RETURNING *;


-- name: GetPostByID :one
SELECT *
FROM posts
WHERE id = $1
LIMIT 1;

-- name: GetPostsByFeedID :many
SELECT *
FROM posts
WHERE feed_id = $1
ORDER BY published_at DESC;

-- name: GetAllPosts :many
SELECT
    id,
    feed_id,
    title,
    url,
    description,
    published_at,
    guid,
    created_at
FROM posts
ORDER BY published_at DESC;
