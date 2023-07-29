CREATE TABLE IF NOT EXISTS services (
    guid uuid NOT NULL,
    group_id uuid NOT NULL,
    name text NOT NULL DEFAULT ''::character varying,
    price numeric NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT services_pkey PRIMARY KEY (guid)
);

ALTER TABLE IF EXISTS services ADD CONSTRAINT "service_service_group_id_fkey" FOREIGN KEY(group_id) REFERENCES service_groups("guid"); 