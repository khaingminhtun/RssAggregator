-- +goose Up
-- ============================
-- USERS
-- ============================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    auth_type TEXT NOT NULL DEFAULT 'local',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Auto-update updated_at
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- ... rest of tables, feeds, posts, indexes ...

-- ============================
-- FEEDS (Global feed metadata)
-- ============================
CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    feed_url TEXT NOT NULL UNIQUE,
    website_url TEXT NOT NULL,
    description TEXT,
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
    guid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT post_feed_guid_unique UNIQUE (feed_id, guid)
);

CREATE INDEX idx_posts_feed_id ON posts(feed_id);

-- Prevent duplicate links per feed
CREATE UNIQUE INDEX post_feed_unique_link
ON posts (feed_id, url);

-- +goose Down
DROP INDEX IF EXISTS post_feed_unique_link;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS user_feeds;
DROP TABLE IF EXISTS feeds;
DROP TRIGGER IF EXISTS users_set_updated_at ON users;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS users;
