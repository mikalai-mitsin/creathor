CREATE TABLE public.sessions
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT sessions_pk PRIMARY KEY,
    title varchar NOT NULL,
    description text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE INDEX search_sessions
    ON public.sessions
        USING GIN (to_tsvector('english', description));

CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE
    ON
        public.sessions
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('session_list', 'Session list'),
       ('session_detail', 'Session detail'),
       ('session_create', 'Session create'),
       ('session_update', 'Session update'),
       ('session_delete', 'Session delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'session_list'),
       ('admin', 'session_detail'),
       ('admin', 'session_create'),
       ('admin', 'session_update'),
       ('admin', 'session_delete'),
       ('user', 'session_list'),
       ('user', 'session_detail'),
       ('user', 'session_create'),
       ('user', 'session_update'),
       ('user', 'session_delete'),
       ('guest', 'session_list'),
       ('guest', 'session_detail');
