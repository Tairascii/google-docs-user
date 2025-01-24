-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id uuid primary key not null default gen_random_uuid(),
    name varchar(64) not null,
    email varchar(255) not null unique,
    password_hash text,
    profile_picture_url text,
    created_at timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
