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
