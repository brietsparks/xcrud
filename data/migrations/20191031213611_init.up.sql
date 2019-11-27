create table "user"
(
    id          serial primary key,
    first_name  varchar(100) not null,
    last_name   varchar(100) not null
);

create table "group"
(
    id          serial primary key,
    name        varchar(100) not null
);

create table group_user
(
    group_id int,
    user_id int,
    primary key (group_id, user_id),
    foreign key (group_id) references "group" (id),
    foreign key (user_id) references "user" (id)
);
