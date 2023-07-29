CREATE TABLE IF NOT EXISTS service_groups (
    guid uuid NOT NULL,
    name text NOT NULL DEFAULT ''::character varying,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT service_groups_pkey PRIMARY KEY (guid)
);
