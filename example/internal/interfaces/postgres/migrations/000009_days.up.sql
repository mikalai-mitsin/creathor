CREATE TABLE public.days
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT days_pk PRIMARY KEY,
    name varchar NOT NULL,
    repeat int NOT NULL,
    equipment_id text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE INDEX search_days
    ON public.days
        USING GIN (to_tsvector('english', name));

CREATE TRIGGER update_days_updated_at
    BEFORE UPDATE
    ON
        public.days
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('day_list', 'Day list'),
       ('day_detail', 'Day detail'),
       ('day_create', 'Day create'),
       ('day_update', 'Day update'),
       ('day_delete', 'Day delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'day_list'),
       ('admin', 'day_detail'),
       ('admin', 'day_create'),
       ('admin', 'day_update'),
       ('admin', 'day_delete'),
       ('user', 'day_list'),
       ('user', 'day_detail'),
       ('user', 'day_create'),
       ('user', 'day_update'),
       ('user', 'day_delete'),
       ('guest', 'day_list'),
       ('guest', 'day_detail');
