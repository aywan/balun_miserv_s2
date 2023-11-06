-- +goose Up
create table "audit" (
    id bigserial primary key,
    created_at timestamp not null default now(),
    creator_id integer,
    reference text not null,
    reference_id int not null,
    action text not null
);

-- +goose Down
drop table "audit";
