-- name: CreateUser :one 
insert into users(name, email,password_hash) 
values ($1, $2,$3)
returning *;

-- name: GetUserById :one
select * from users where id = $1;

-- name: GetUserByEmail :one 
select * from users where email = $1;