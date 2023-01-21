create table public.permissions
(
    id   varchar(255) not null
        constraint permissions_pk
            primary key,
    name varchar(255)
);
