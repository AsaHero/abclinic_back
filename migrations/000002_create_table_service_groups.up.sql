CREATE TABLE IF NOT EXISTS service_groups (
    guid uuid NOT NULL,
    name text NOT NULL DEFAULT ''::character varying,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT service_groups_pkey PRIMARY KEY (guid)
);

ALTER TABLE IF EXISTS services ADD CONSTRAINT "service_groups_group_id_fkey" FOREIGN KEY(group_id) REFERENCES service_groups("guid"); 