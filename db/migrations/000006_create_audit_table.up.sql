create table if not exists audits
(
    id varchar(40) primary key,
    task_id varchar(40) references tasks (id),
    comment varchar(100) unique not null,
    status varchar(50),
    created_by varchar(40) references users (id),
    updated_by varchar(40) references users (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
