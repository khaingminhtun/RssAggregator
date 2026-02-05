-- name: CreateSubscription :one
INSERT INTO user_feeds(user_id, feed_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserSubscribedFeeds :many
SELECT
    f.id,
    f.title,
    f.feed_url,
    f.website_url,
    f.description,
    f.last_fetched_at,
    f.created_at
FROM user_feeds uf
JOIN feeds f ON uf.feed_id = f.id
WHERE uf.user_id = $1
ORDER BY uf.created_at DESC;

-- name: UnsubscribeUserFromFeed :exec
DELETE FROM user_feeds
WHERE user_id = $1 AND feed_id = $2;
