-- +goose Up
-- ============================
-- Make feeds.site_url NOT NULL
-- ============================

-- Step 1: Alter the column to NOT NULL
ALTER TABLE feeds
ALTER COLUMN site_url SET NOT NULL;

-- +goose Down
-- ============================
-- Revert feeds.site_url to nullable
-- ============================
ALTER TABLE feeds
ALTER COLUMN site_url DROP NOT NULL;
