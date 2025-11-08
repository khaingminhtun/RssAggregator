-- +goose Up

-- ---------------------------
-- 1. user_post_status Table (Tracking user interaction with posts)
-- ---------------------------
CREATE TABLE user_post_status (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    is_read boolean NOT NULL DEFAULT FALSE,
    is_favorite boolean NOT NULL DEFAULT FALSE,
    
    -- Composite Primary Key: Ensures one status record per user/post pair
    PRIMARY KEY (user_id, post_id)
);


-- +goose Down
DROP TABLE user_post_status;
