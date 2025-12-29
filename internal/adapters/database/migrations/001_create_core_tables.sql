-- +goose Up
-- ============================
-- USERS (Authentication)
-- ============================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================
-- FEEDS (Global feed metadata)
-- ============================
CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    last_fetched_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================
-- USER_FEEDS (Subscriptions)
-- ============================
CREATE TABLE user_feeds (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT user_feed_uniqueness UNIQUE (user_id, feed_id)
);

-- Indexes for FK lookups
CREATE INDEX idx_user_feeds_user_id ON user_feeds(user_id);
CREATE INDEX idx_user_feeds_feed_id ON user_feeds(feed_id);

-- ============================
-- POSTS (Aggregated content)
-- ============================
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMPTZ NOT NULL,
    guid TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Prevent duplicate items per feed
    CONSTRAINT post_feed_guid_unique UNIQUE (feed_id, guid)
);

-- Index for feed → posts queries
CREATE INDEX idx_posts_feed_id ON posts(feed_id);

-- +goose Down
DROP TABLE posts;
DROP TABLE user_feeds;
DROP TABLE feeds;
DROP TABLE users;
