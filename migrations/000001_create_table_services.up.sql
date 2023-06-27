CREATE TABLE IF NOT EXISTS services (
    guid uuid NOT NULL,
    group_id uuid NOT NULL,
    name text NOT NULL DEFAULT''::character varying,
    price numeric NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT price_list_pkey PRIMARY KEY (guid)
);