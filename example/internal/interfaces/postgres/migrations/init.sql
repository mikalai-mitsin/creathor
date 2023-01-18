CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create function update_updated_at_task() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

create table public.permissions
(
    id   varchar(255) not null
        constraint permissions_pk
            primary key,
    name varchar(255)
);

create table public.groups
(
    id   varchar(255) not null
        constraint groups_pk
            primary key,
    name varchar(255) not null
);

create table public.group_permissions
(
    group_id      varchar(255) not null
        references public.groups
            on update cascade on delete cascade
        constraint group_permissions_group_ids
            references public.groups
            deferrable initially deferred,
    permission_id varchar(255) not null
        references public.permissions
            on update cascade on delete cascade
        constraint group_permissions_permission_ids
            references public.permissions
            deferrable initially deferred
);

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
