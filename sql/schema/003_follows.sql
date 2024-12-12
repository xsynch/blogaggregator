-- +goose Up
CREATE TABLE feed_follow (id uuid primary key, created_at timestamp not null, updated_at timestamp not null,  user_id uuid not null references users(id) on delete cascade, feed_id uuid not null references feeds(id) on delete cascade, constraint user_feed unique(user_id, feed_id));

-- +goose Down
DROP TABLE feed_follow;