-- name : CreateeUser : one user
insert into users(name, email,password_hash) 
values ($1, $2,$3)
returning *;

--name : GetUserById : one user by id
select * from users where id = $1;

-- name : GetUserByEmail : one user by email
select * from users where email = $1;