-- +goose Up
ALTER TABLE users
ADD COLUMN auth_type TEXT NOT NULL DEFAULT 'local';

-- +goose Down
ALTER TABLE users
DROP COLUMN auth_type;
