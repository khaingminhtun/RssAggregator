-- +goose Up
-- ============================
-- USER_POSTS (Per-user post state)
-- ============================
CREATE TABLE user_posts (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,

    is_read BOOLEAN NOT NULL DEFAULT false,
    is_favorite BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, post_id)
);

-- Auto-update updated_at
-- +goose StatementBegin
CREATE TRIGGER user_posts_set_updated_at
BEFORE UPDATE ON user_posts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

CREATE INDEX idx_user_posts_user_id ON user_posts(user_id);
CREATE INDEX idx_user_posts_post_id ON user_posts(post_id);
CREATE INDEX idx_user_posts_favorite ON user_posts(user_id, is_favorite);

-- +goose Down
DROP TRIGGER IF EXISTS user_posts_set_updated_at ON user_posts;
DROP TABLE IF EXISTS user_posts;
