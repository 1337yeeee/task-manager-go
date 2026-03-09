create table if not exists users
(
    id varchar(40) primary key,
    name varchar(50) not null,
    email varchar(255) unique not null,
    password varchar(255) not null,
    role varchar(255) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
