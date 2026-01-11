-- +goose Up
-- rename url to feed_url and site_url to webstie_url
ALTER TABLE feeds
    RENAME COLUMN url TO feed_url;

ALTER TABLE feeds
    RENAME COLUMN site_url TO website_url;

-- +goose Down
-- rollback column names
ALTER TABLE feeds
    RENAME COLUMN feed_url TO url;

ALTER TABLE feeds
    RENAME COLUMN website_url TO site_url;


