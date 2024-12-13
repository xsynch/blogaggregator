-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
select * from users where name = $1;

-- name: GetUserByid :one
select * from users where id = $1;


-- name: Reset :exec
delete from users;

-- name: GetUsers :many
select * from users;

-- name: CreateFeed :one
insert into feeds (id,created_at,updated_at,name,url,user_id)
values ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetFeedsByID :one
select * from feeds where id = $1;

-- name: CreateFeedFollow :one 
with inserted_feed_follow as (insert into feed_follow(id,created_at,updated_at,user_id,feed_id)
values ($1,$2,$3,$4,$5)
RETURNING *) select inserted_feed_follow.*,
feeds.name as feed_name,
users.name as user_name
from inserted_feed_follow
inner join feeds on inserted_feed_follow.feed_id = feeds.id
inner join users on inserted_feed_follow.user_id = users.id;

-- name: GetFeedIDByName :one
select * from feeds where url = $1;

-- name: GetFeedFollowsForUser :many
select * from feed_follow where user_id = $1;

-- name: Unfollow :exec
delete from feed_follow where feed_id = $1 and user_id = $2;

-- name: MarkFeedFetched :exec
update feeds set last_fetched_at = $1, updated_at=$1 where id = $2;

-- name: GetNextFeedToFetch :one
select * from feeds order by last_fetched_at asc nulls first limit 1;

-- name: CreatePost :one
insert into posts (id , created_at , updated_at , title , url , description , published_at ,feed_id ) values($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- name: GetPostsForUser :many
select * from posts where feed_id = $1 limit $2;