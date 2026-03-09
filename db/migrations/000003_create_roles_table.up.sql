create table if not exists roles
(
    id varchar(40) primary key,
    user_id varchar(40) references users (id),
    role varchar(255) unique not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
