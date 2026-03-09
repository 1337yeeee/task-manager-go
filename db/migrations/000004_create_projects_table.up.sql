create table if not exists projects
(
    id varchar(40) primary key,
    name varchar(100) unique not null,
    description text,
    content text,
    created_by varchar(40) references users (id),
    updated_by varchar(40) references users (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
