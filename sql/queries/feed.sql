-- name: CreateFeed :one
insert into feeds(title, url)
values ($1, $2)
on conflict (url) do nothing
returning *;


-- name: GetFeedByUrl :one
select * from feeds where url = $1;

-- name: GetFeedsToFetch :many
select * from feeds
where last_fetched_at < now() - interval '30 minutes'
order by last_fetched_at asc
limit $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = now()
where id = $1;

-- name: ListFeeds :many
select * from feeds
order by created_at desc;