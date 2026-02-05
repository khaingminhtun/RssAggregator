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


-- name: GetLatestPosts :many
SELECT *
FROM posts
ORDER BY published_at DESC
LIMIT $1;


-- name: SearchPosts :many
SELECT *
FROM posts
WHERE title ILIKE '%' || $1 || '%'
   OR description ILIKE '%' || $1 || '%'
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPostsForUser :many
SELECT
    p.id,
    p.feed_id,
    p.title,
    p.url,
    p.description,
    p.published_at,
    p.guid,
    p.created_at
FROM posts p
JOIN user_feeds uf ON p.feed_id = uf.feed_id
WHERE uf.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2 OFFSET $3;
