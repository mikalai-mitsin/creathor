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
