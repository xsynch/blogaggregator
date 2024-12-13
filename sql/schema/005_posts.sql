-- +goose Up
CREATE TABLE posts (id uuid primary key, created_at timestamp not null, updated_at timestamp not null, title text not null, url text unique not null, description text, published_at timestamp, feed_id uuid not null references feeds(id) on delete cascade);

-- +goose Down
DROP TABLE posts;