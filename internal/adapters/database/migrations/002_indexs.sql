-- +goose Up
-- ============================
-- Fix posts table and add performance indexes
-- ============================

-- ----------------------------
-- POSTS: fix GUID and remove duplicate URL index
-- ----------------------------

-- Drop the old duplicate URL index if it exists
DROP INDEX IF EXISTS post_feed_unique_link;

-- Ensure guid is NOT NULL (for proper upsert)
ALTER TABLE posts
ALTER COLUMN guid SET NOT NULL;

-- ----------------------------
-- POSTS: performance indexes
-- ----------------------------

-- Index for fetching posts by feed in published order
CREATE INDEX IF NOT EXISTS idx_posts_feed_published_at
ON posts(feed_id, published_at DESC);

-- Index for global latest posts
CREATE INDEX IF NOT EXISTS idx_posts_published_at
ON posts(published_at DESC);

-- ----------------------------
-- USER_FEEDS: performance indexes
-- ----------------------------

-- Index for joining user_feeds to posts
CREATE INDEX IF NOT EXISTS idx_user_feeds_user_feed
ON user_feeds(user_id, feed_id);

-- Optional: full-text search (enable pg_trgm extension first)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_posts_title_trgm
ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_desc_trgm
ON posts USING gin (description gin_trgm_ops);

-- +goose Down
-- ============================
-- Drop the indexes added above and revert guid column if needed
-- ============================

-- Drop performance indexes
DROP INDEX IF EXISTS idx_posts_feed_published_at;
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_user_feeds_user_feed;

-- Optional: drop pg_trgm search indexes
DROP INDEX IF EXISTS idx_posts_title_trgm;
DROP INDEX IF EXISTS idx_posts_desc_trgm;

-- Revert guid to nullable (if needed)
ALTER TABLE posts
ALTER COLUMN guid DROP NOT NULL;
