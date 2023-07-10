CREATE TABLE IF NOT EXISTS chapters (
    guid uuid NOT NULL,
    title character varying(256) DEFAULT ''::character varying,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT chapters_pkey PRIMARY KEY (guid)
);