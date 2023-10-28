-- +goose Up
create table "user" (
    id serial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp null,
    name text not null,
    email text not null,
    password_hash text not null,
    role int not null default 0
);

create unique index "user_email_udx" on "user"(email) where deleted_at is null;

-- +goose Down
drop table "user";
