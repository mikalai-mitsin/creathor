CREATE TABLE public.sessions
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT sessions_pk PRIMARY KEY,
    title varchar NOT NULL,
    description text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE
    ON
        public.sessions
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();
