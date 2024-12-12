-- +goose Up
CREATE TABLE feeds (id uuid primary key, created_at timestamp not null, updated_at timestamp not null, name text unique not null, url text unique not null, user_id uuid not null constraint fk_users references users(id) on delete cascade);

-- +goose Down
DROP TABLE feeds;