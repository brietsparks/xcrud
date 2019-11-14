create table users
(
    id          serial primary key,
    first_name  varchar(100) not null,
    last_name   varchar(100) not null
);

create table groups
(
    id          serial primary key,
    name        varchar(100) not null
);

create table groups_users
(
    group_id int,
    user_id int,
    primary key (group_id, user_id),
    foreign key (group_id) references groups (id),
    foreign key (user_id) references users (id)
);
