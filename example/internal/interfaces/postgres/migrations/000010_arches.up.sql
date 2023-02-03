CREATE TABLE public.arches
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT arches_pk PRIMARY KEY,
    name varchar NOT NULL,
    tags varchar[] NOT NULL,
    versions bigint[] NOT NULL,
    old_versions bigint[] NOT NULL,
    release timestamp NOT NULL,
    tested timestamp NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE INDEX search_arches
    ON public.arches
        USING GIN (to_tsvector('english', name));

CREATE TRIGGER update_arches_updated_at
    BEFORE UPDATE
    ON
        public.arches
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('arch_list', 'Arch list'),
       ('arch_detail', 'Arch detail'),
       ('arch_create', 'Arch create'),
       ('arch_update', 'Arch update'),
       ('arch_delete', 'Arch delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'arch_list'),
       ('admin', 'arch_detail'),
       ('admin', 'arch_create'),
       ('admin', 'arch_update'),
       ('admin', 'arch_delete'),
       ('user', 'arch_list'),
       ('user', 'arch_detail'),
       ('user', 'arch_create'),
       ('user', 'arch_update'),
       ('user', 'arch_delete'),
       ('guest', 'arch_list'),
       ('guest', 'arch_detail');
