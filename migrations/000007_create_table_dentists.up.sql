CREATE TABLE IF NOT EXISTS dentists (
    id serial,
    clone_name character varying NOT NULL,
    name character varying,
    info text,
    url character varying,
    side character varying NOT NULL,
    priority numeric NOT NULL
);