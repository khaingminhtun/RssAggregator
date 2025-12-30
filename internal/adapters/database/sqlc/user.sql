-- name: CreateUser :one 
INSERT INTO users (name, email, password_hash, auth_type)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
