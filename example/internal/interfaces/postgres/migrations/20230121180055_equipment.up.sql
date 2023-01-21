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
        USING GIN (to_tsvector('english', title || body));

CREATE TRIGGER update_equipment_updated_at
    BEFORE UPDATE
    ON
        public.equipment
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();
