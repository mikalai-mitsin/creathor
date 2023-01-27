CREATE TABLE public.plans
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT plans_pk PRIMARY KEY,
    name varchar NOT NULL,
    repeat bigint NOT NULL,
    equipment_id text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE INDEX search_plans
    ON public.plans
        USING GIN (to_tsvector('english', name));

CREATE TRIGGER update_plans_updated_at
    BEFORE UPDATE
    ON
        public.plans
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('plan_list', 'Plan list'),
       ('plan_detail', 'Plan detail'),
       ('plan_create', 'Plan create'),
       ('plan_update', 'Plan update'),
       ('plan_delete', 'Plan delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'plan_list'),
       ('admin', 'plan_detail'),
       ('admin', 'plan_create'),
       ('admin', 'plan_update'),
       ('admin', 'plan_delete'),
       ('user', 'plan_list'),
       ('user', 'plan_detail'),
       ('user', 'plan_create'),
       ('user', 'plan_update'),
       ('user', 'plan_delete'),
       ('guest', 'plan_list'),
       ('guest', 'plan_detail');
