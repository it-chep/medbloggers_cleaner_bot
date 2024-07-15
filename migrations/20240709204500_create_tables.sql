-- +goose Up
-- +goose StatementBegin
create table if not exists tg_users
(
    id                     bigserial primary key,
    tg_id                  bigint,
    name                   varchar(255),
    username               varchar(255) null,
    last_message_timestamp timestamp    null,
    message_counter        int     default 1,
    chat_id                bigint,
    admin                  boolean default false,
    is_blocked             boolean default false
);

create table if not exists chats
(
    id      bigserial primary key,
    chat_id bigint,
    name    text
);

create table if not exists hashtags
(
    id      bigserial primary key,
    chat_id bigint,
    hashtag text
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tg_users;
drop table if exists chats;
drop table if exists hashtags
-- +goose StatementEnd
