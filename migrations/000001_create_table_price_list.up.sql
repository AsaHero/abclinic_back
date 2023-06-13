CREATE TABLE IF NOT EXISTS price_list (
    guid uuid NOT NULL,
    title character varying(256) NOT NULL DEFAULT ''::character varying,
    service text NOT NULL DEFAULT ''::character varying,
    price numeric NOT NULL,
    language character varying(2),
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT price_list_pkey PRIMARY KEY (guid)
);