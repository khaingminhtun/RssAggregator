-- +goose Up
ALTER TABLE users ALTER COLUMN auth_type SET DEFAULT 'local';
UPDATE users SET auth_type = 'local';

-- +goose Down
-- (optional)
