create table public.users
(
    id         uuid         default uuid_generate_v4()               not null
        constraint users_pk
            primary key,
    first_name varchar(255)                                          not null,
    last_name  varchar(255)                                          not null,
    password   varchar(255)                                          not null,
    email      varchar(255)                                          not null
        unique,
    group_id   varchar(255) default 'user'::character varying        not null
        constraint group_id_fk
            references public.groups
            on update cascade on delete cascade,
    created_at timestamp    default (now() AT TIME ZONE 'utc'::text) not null,
    updated_at timestamp    default (now() AT TIME ZONE 'utc'::text) not null
);

create index search_users
    on public.users using gin (to_tsvector('english'::regconfig, (first_name::text || last_name::text) || email::text));

create trigger update_users_updated_at
    before update
    on public.users
    for each row
execute procedure public.update_updated_at_task();

insert into public.permissions (id, name)
values ('user_list', 'User list'),
       ('user_detail', 'User detail'),
       ('user_create', 'User create'),
       ('user_update', 'User update'),
       ('user_delete', 'User delete');

insert into public.group_permissions (group_id, permission_id)
values ('admin', 'user_list'),
       ('admin', 'user_detail'),
       ('admin', 'user_create'),
       ('admin', 'user_update'),
       ('admin', 'user_delete'),
       ('user', 'user_list'),
       ('user', 'user_detail'),
       ('user', 'user_create'),
       ('user', 'user_update'),
       ('user', 'user_delete'),
       ('guest', 'user_list'),
       ('guest', 'user_detail'),
       ('guest', 'user_create');
