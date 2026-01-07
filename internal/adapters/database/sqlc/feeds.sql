-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: CreateFeed :one
INSERT INTO feeds (url, title, description, site_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeedURLByID :one
SELECT url
FROM feeds
WHERE id = $1;
