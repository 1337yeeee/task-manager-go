create table if not exists auths
(
    id varchar(40) primary key,
    user_id varchar(40) references users (id),
    access_token TEXT unique not null,
    refresh_token TEXT unique not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    expires_at timestamp
);
