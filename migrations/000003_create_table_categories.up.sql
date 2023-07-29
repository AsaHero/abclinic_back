CREATE TABLE IF NOT EXISTS categories (
    guid uuid NOT NULL,
    title character varying(256) DEFAULT ''::character varying,
    description character varying(512) DEFAULT ''::character varying,
    url character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT categories_pkey PRIMARY KEY (guid)
);