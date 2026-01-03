-- name: CreateSubscription :one
INSERT INTO user_feeds(user_id, feed_id)
VALUES ($1, $2)
RETURNING *;
