CREATE TABLE IF NOT EXISTS authors (
    guid uuid NOT NULL,
    name character varying,
    img bytea,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT authors_pkey PRIMARY KEY (guid)
);