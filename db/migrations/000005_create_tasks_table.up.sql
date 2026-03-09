create table if not exists tasks
(
    id varchar(40) primary key,
    project_id varchar(40) references projects (id),
    name varchar(100) unique not null,
    content text,
    status varchar(50),
    executive_id varchar(40) references users (id) not null,
    auditor_id varchar(40) references users (id) null,
    created_by varchar(40) references users (id),
    updated_by varchar(40) references users (id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
