CREATE TABLE public.equipment
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT equipment_pk PRIMARY KEY,
    name varchar NOT NULL,
    repeat int NOT NULL,
    weight int NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE INDEX search_equipment
    ON public.equipment
        USING GIN (to_tsvector('english', name));

CREATE TRIGGER update_equipment_updated_at
    BEFORE UPDATE
    ON
        public.equipment
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('equipment_list', 'Equipment list'),
       ('equipment_detail', 'Equipment detail'),
       ('equipment_create', 'Equipment create'),
       ('equipment_update', 'Equipment update'),
       ('equipment_delete', 'Equipment delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'equipment_list'),
       ('admin', 'equipment_detail'),
       ('admin', 'equipment_create'),
       ('admin', 'equipment_update'),
       ('admin', 'equipment_delete'),
       ('user', 'equipment_list'),
       ('user', 'equipment_detail'),
       ('user', 'equipment_create'),
       ('user', 'equipment_update'),
       ('user', 'equipment_delete'),
       ('guest', 'equipment_list'),
       ('guest', 'equipment_detail');
