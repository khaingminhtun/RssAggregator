-- +goose Up
UPDATE users SET auth_type = 'local' WHERE auth_type IS NULL;

-- +goose Down
-- (optional)
