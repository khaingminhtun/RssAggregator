-- +goose Up
ALTER TABLE users ALTER COLUMN auth_type DROP NOT NULL;
ALTER TABLE users ALTER COLUMN auth_type SET DEFAULT 'local';
UPDATE users SET auth_type = 'local' WHERE auth_type IS NULL;

-- +goose Down
ALTER TABLE users ALTER COLUMN auth_type SET NOT NULL;
ALTER TABLE users ALTER COLUMN auth_type DROP DEFAULT;
