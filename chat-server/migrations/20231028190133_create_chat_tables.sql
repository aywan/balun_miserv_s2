-- +goose Up
create table chat (
    id serial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp null,
    owner_id int,
    name text not null
);

create table chat_user (
    chat_id int not null,
    user_id int not null,
    last_message_id bigint,
    foreign key (chat_id) references chat (id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table message (
    id bigserial primary key,
    created_at timestamp not null default now(),
    chat_id int not null,
    user_id int,
    type int not null default 0,
    text text not null,
    foreign key (chat_id) references chat (id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table chat_message (
    chat_id int not null,
    last_message_id bigint not null,
    foreign key (chat_id) references chat (id) ON UPDATE CASCADE ON DELETE CASCADE,
    foreign key (last_message_id) references message (id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- +goose Down
drop table chat_message;
drop table message;
drop table chat_user;
drop table chat;
