-- name: CreateFeed :one
INSERT INTO feeds (feed_url, title, description, website_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeedURLByID :one
SELECT feed_url
FROM feeds
WHERE id = $1;

-- name: GetFeedByID :one
SELECT *
FROM feeds
WHERE id = $1;

-- name: GetAllFeeds :many
SELECT *
FROM feeds
ORDER BY created_at DESC;

-- name: GetFeedsByUserID :many
SELECT f.*
FROM feeds f
JOIN user_feeds uf ON uf.feed_id = f.id
WHERE uf.user_id = $1
ORDER BY uf.created_at DESC;

-- name: UpdateFeed :one
UPDATE feeds
SET title = $2,
    feed_url = $3,
    description = $4,
    website_url = $5,
    last_fetched_at = $6
WHERE id = $1
RETURNING *;

-- name: UpdateFeedFetchTime :one
UPDATE feeds
SET last_fetched_at = $2
WHERE id = $1
RETURNING *;


-- name: DeleteFeedByID :exec
DELETE FROM feeds
WHERE id = $1;

-- name: DeleteFeedIfUnused :exec
DELETE FROM feeds f
WHERE f.id = $1
  AND NOT EXISTS (
      SELECT 1
      FROM user_feeds uf
      WHERE uf.feed_id = f.id
  );

-- name: DeleteAllFeedsByUserID :exec
DELETE FROM feeds f
USING user_feeds uf
WHERE uf.feed_id = f.id
  AND uf.user_id = $1
  AND NOT EXISTS (
      SELECT 1
      FROM user_feeds uf2
      WHERE uf2.feed_id = f.id
        AND uf2.user_id <> $1
  );
