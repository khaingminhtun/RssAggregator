-- +goose Up
-- ============================
-- Extend feeds metadata
-- ============================
ALTER TABLE feeds
ADD COLUMN site_url TEXT,
ADD COLUMN description TEXT;

-- ============================
-- Posts improvements
-- ============================
ALTER TABLE posts
ALTER COLUMN guid DROP NOT NULL;

CREATE UNIQUE INDEX post_feed_unique_link
ON posts (feed_id, url);

-- +goose Down
DROP INDEX post_feed_unique_link;

ALTER TABLE posts
ALTER COLUMN guid SET NOT NULL;

ALTER TABLE feeds
DROP COLUMN description,
DROP COLUMN site_url;
