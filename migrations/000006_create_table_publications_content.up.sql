CREATE TABLE IF NOT EXISTS publications_content (
    guid uuid NOT NULL,
    publication_id uuid NOT NULL,
    order_id numeric DEFAULT 0,
    url character varying,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT publications_content_pkey PRIMARY KEY (guid)
);

ALTER TABLE IF EXISTS publications_content ADD CONSTRAINT "publications_content_publication_id_fkey" FOREIGN KEY(publication_id) REFERENCES publications("guid"); 