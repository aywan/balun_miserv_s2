-- +goose Up

alter table chat_message add constraint chat_message_chat_id_pk primary key (chat_id);
alter table chat_user add constraint chat_user_chat_id_user_id_pk primary key (chat_id, user_id);

-- +goose Down
